# xk6-coap

`xk6-coap` is a [`k6` extension](https://k6.io/docs/extensions/) for the
[Constrained Application Protocol
(CoAP)](https://www.rfc-editor.org/rfc/rfc7252).

## Getting Started

Using `k6` extentions requires including the extension(s) in a `k6` build. The
[`xk6`](https://github.com/grafana/xk6) tool will handle executing the build,
and the `Makefile` in this repository will ensure that `xk6` is installed and
produces a valid build.

```
make build
```

This will produce a `k6` binary in the current working directory. It can be used
to execute a simple load test (`simple.js`) that will spin up the specified
number of virtual users (vu's) and establish observations for 10 seconds.
`xk6-coap` currently requires the use of Pre-Shared Keys (PSKs) for connecting
to CoAP endpoints. By default, the values of the `COAP_PSK` and `COAP_PSK_ID`
environment variables will be used for the key and ID respectively.

To execute a test with 2 virtual users making connections and sending `GET`
messages for 5 seconds, the following command could be run.

```
./k6 run ./examples/simple.js --vus 10 --duration 5s
```

Reference the [`k6` documentation](https://k6.io/docs/using-k6/test-lifecycle/)
for more information on how to configure and run tests.
