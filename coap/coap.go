package coap

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dop251/goja"
	"github.com/mstoykov/k6-taskqueue-lib/taskqueue"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	udp "github.com/plgd-dev/go-coap/v3/udp/client"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

const (
	endpointArgIdx = 1
	pskIDEnvArgIdx = 2
	pskEnvArgIdx   = 3
	certPathArgIdx = 4
	keyPathArgIdx  = 5
)

// Message is a CoAP message.
type Message struct {
	Code string
	Body []byte
}

// RootModule is the imported module for tests using xk6-coap.
type RootModule struct{}

// CoAP constructs new Constrained Application Protocol clients.
type CoAP struct {
	vu modules.VU
}

var _ modules.Module = &RootModule{}

// NewModuleInstance instantiates a new root xk6-coap module.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &CoAP{
		vu: vu,
	}
}

// Exports defines the ESM exports for a CoAP instance.
func (c *CoAP) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Client": c.client,
		},
	}
}

// client constructs a new CoAP client by establishing a new DTLS cnnection with
// the provided server endpoint.
func (c *CoAP) client(cc goja.ConstructorCall) *goja.Object {
	rt := c.vu.Runtime()
	endpoint := cc.Argument(endpointArgIdx).String()
	conf := &piondtls.Config{
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(context.Background(), 10*time.Second)
		},
	}

	// Only ECDSA keys are currently supported.
	if !goja.IsUndefined(cc.Argument(certPathArgIdx)) && !goja.IsUndefined(cc.Argument(keyPathArgIdx)) {
		pemCert, err := os.ReadFile(filepath.Clean(cc.Argument(certPathArgIdx).String()))
		if err != nil {
			common.Throw(rt, err)
			return nil
		}
		certBlock, _ := pem.Decode(pemCert)
		pemKey, err := os.ReadFile(filepath.Clean(cc.Argument(certPathArgIdx).String()))
		if err != nil {
			common.Throw(rt, err)
			return nil
		}
		keyBlock, _ := pem.Decode(pemKey)
		key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			common.Throw(rt, err)
			return nil
		}
		conf.Certificates = []tls.Certificate{
			{
				Certificate: [][]byte{certBlock.Bytes},
				PrivateKey:  key,
			},
		}
		conf.CipherSuites = []piondtls.CipherSuiteID{
			piondtls.TLS_ECDHE_ECDSA_WITH_AES_128_CCM,
			piondtls.TLS_ECDHE_ECDSA_WITH_AES_128_CCM_8,
			piondtls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		}
	}

	// If certificates were provided, they take precedence over PSK as
	// piondtls will always use PSK if provided.
	if len(conf.Certificates) == 0 && !goja.IsUndefined(cc.Argument(pskIDEnvArgIdx)) && !goja.IsUndefined(cc.Argument(pskEnvArgIdx)) {
		pskID, _ := os.LookupEnv(cc.Argument(pskIDEnvArgIdx).String())
		conf.PSKIdentityHint = []byte(pskID)
		psk, _ := os.LookupEnv(cc.Argument(pskEnvArgIdx).String())
		conf.PSK = func(hint []byte) ([]byte, error) {
			return []byte(psk), nil
		}
		conf.CipherSuites = []piondtls.CipherSuiteID{
			piondtls.TLS_PSK_WITH_AES_128_CBC_SHA256,
			piondtls.TLS_PSK_WITH_AES_128_GCM_SHA256,
			piondtls.TLS_PSK_WITH_AES_128_CCM_8,
			piondtls.TLS_PSK_WITH_AES_128_CCM,
		}
	}

	conn, err := dtls.Dial(endpoint, conf)
	if err != nil {
		common.Throw(rt, err)
		return nil
	}

	client := &client{
		vu:   c.vu,
		tq:   taskqueue.New(c.vu.RegisterCallback),
		conn: conn,
		obj:  rt.NewObject(),
	}

	if err := client.obj.DefineDataProperty("get", rt.ToValue(client.Get), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	if err := client.obj.DefineDataProperty("observe", rt.ToValue(client.Observe), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	if err := client.obj.DefineDataProperty("put", rt.ToValue(client.Put), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	if err := client.obj.DefineDataProperty("post", rt.ToValue(client.Post), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	if err := client.obj.DefineDataProperty("delete", rt.ToValue(client.Delete), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	if err := client.obj.DefineDataProperty("close", rt.ToValue(client.Close), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
		return nil
	}
	return client.obj
}

// client is a CoAP client with a DTLS connection.
type client struct {
	vu   modules.VU
	tq   *taskqueue.TaskQueue
	conn *udp.Conn
	obj  *goja.Object
}

// Get sends a GET message to the specified path.
func (c *client) Get(path string, timeout int) Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	msg, err := c.conn.Get(ctx, path)
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
		return Message{}
	}
	var b []byte
	if body := msg.Body(); body != nil {
		if b, err = io.ReadAll(body); err != nil {
			common.Throw(c.vu.Runtime(), err)
			return Message{}
		}
	}
	return Message{
		Code: msg.Code().String(),
		Body: b,
	}
}

// Observe sends an OBSERVE message to the specified path. It waits for messages
// until the specified timeout.
func (c *client) Observe(path string, timeout int, listener func(Message)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))

	obs, err := c.conn.Observe(ctx, path, func(msg *pool.Message) {
		var b []byte
		var err error
		if body := msg.Body(); body != nil {
			if b, err = io.ReadAll(body); err != nil {
				common.Throw(c.vu.Runtime(), err)
				return
			}
		}
		c.tq.Queue(func() error {
			listener(Message{
				Code: msg.Code().String(),
				Body: b,
			})
			return nil
		})
	})
	if err != nil {
		defer cancel()
		common.Throw(c.vu.Runtime(), err)
		return
	}
	go func() {
		<-ctx.Done()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := obs.Cancel(ctx); err != nil {
			c.vu.State().Logger.Warnf("failed to cancel observation: %v", err)
		}
	}()
}

// Put sends a PUT message with the provided content to the specified path.
func (c *client) Put(path, mediaType string, content []byte, timeout int) Message {
	mt, err := message.ToMediaType(mediaType)
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	msg, err := c.conn.Put(ctx, path, mt, bytes.NewReader(content))
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
		return Message{}
	}
	var b []byte
	if body := msg.Body(); body != nil {
		if b, err = io.ReadAll(body); err != nil {
			common.Throw(c.vu.Runtime(), err)
			return Message{}
		}
	}
	return Message{
		Code: msg.Code().String(),
		Body: b,
	}
}

// Post sends a POST message with the provided content to the specified path.
func (c *client) Post(path, mediaType string, content []byte, timeout int) Message {
	mt, err := message.ToMediaType(mediaType)
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	msg, err := c.conn.Post(ctx, path, mt, bytes.NewReader(content))
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
		return Message{}
	}
	var b []byte
	if body := msg.Body(); body != nil {
		if b, err = io.ReadAll(body); err != nil {
			common.Throw(c.vu.Runtime(), err)
			return Message{}
		}
	}
	return Message{
		Code: msg.Code().String(),
		Body: b,
	}
}

// Post sends a POST message with the provided content to the specified path.
func (c *client) Delete(path string, timeout int) Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	msg, err := c.conn.Delete(ctx, path)
	if err != nil {
		common.Throw(c.vu.Runtime(), err)
		return Message{}
	}
	var b []byte
	if body := msg.Body(); body != nil {
		if b, err = io.ReadAll(body); err != nil {
			common.Throw(c.vu.Runtime(), err)
			return Message{}
		}
	}
	return Message{
		Code: msg.Code().String(),
		Body: b,
	}
}

// Close closes the underlying connection.
func (c *client) Close() {
	defer c.tq.Close()
	if err := c.conn.Close(); err != nil {
		common.Throw(c.vu.Runtime(), err)
	}
}
