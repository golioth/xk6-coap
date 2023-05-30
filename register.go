package xk6coap

import (
	"github.com/golioth/xk6-coap/coap"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/coap", new(coap.RootModule))
}
