package coap

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/dop251/goja"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	udp "github.com/plgd-dev/go-coap/v3/udp/client"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

const (
	defaultPSKEnv   = "COAP_PSK"
	defaultPSKIDEnv = "COAP_PSK_ID"
)

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
	endpoint := cc.Argument(0).String()
	pskIDEnv := defaultPSKIDEnv
	if !goja.IsUndefined(cc.Argument(1)) {
		pskIDEnv = cc.Argument(1).String()
	}
	pskID, _ := os.LookupEnv(pskIDEnv)
	pskEnv := defaultPSKEnv
	if !goja.IsUndefined(cc.Argument(2)) {
		pskEnv = cc.Argument(2).String()
	}
	psk, _ := os.LookupEnv(pskEnv)
	conn, err := dtls.Dial(endpoint, &piondtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			return []byte(psk), nil
		},
		PSKIdentityHint: []byte(pskID),
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(context.Background(), 10*time.Second)
		},
		CipherSuites: []piondtls.CipherSuiteID{
			piondtls.TLS_PSK_WITH_AES_128_CCM,
			piondtls.TLS_PSK_WITH_AES_128_CCM_8,
			piondtls.TLS_PSK_WITH_AES_128_GCM_SHA256,
		},
	})
	if err != nil {
		common.Throw(rt, err)
	}

	client := &client{
		conn: conn,
		obj:  rt.NewObject(),
	}

	if err := client.obj.DefineDataProperty("get", rt.ToValue(client.Get), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
	}
	if err := client.obj.DefineDataProperty("observe", rt.ToValue(client.Observe), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
	}
	if err := client.obj.DefineDataProperty("close", rt.ToValue(client.Close), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		common.Throw(rt, err)
	}
	return client.obj
}

// client is a CoAP client with a DTLS connection.
type client struct {
	conn *udp.Conn
	obj  *goja.Object
}

// Get sends a GET message to the specified path.
func (c *client) Get(path string, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	resp, err := c.conn.Get(ctx, path)
	if err != nil {
		return err
	}
	if _, err := io.ReadAll(resp.Body()); err != nil {
		return err
	}
	return nil
}

// Observe sends an OBSERVE message to the specified path. It waits for messages
// until the specified timeout.
func (c *client) Observe(path string, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	obs, err := c.conn.Observe(ctx, path, func(req *pool.Message) {
		// TODO(hasheddan): emit metrics on observed messages.
	})
	if err != nil {
		return err
	}
	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return obs.Cancel(ctx)
}

// Close closes the underlying connection.
func (c *client) Close() error {
	return c.conn.Close()
}
