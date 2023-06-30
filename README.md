# xk6-coap

`xk6-coap` is a [`k6` extension](https://k6.io/docs/extensions/) for the
[Constrained Application Protocol
(CoAP)](https://www.rfc-editor.org/rfc/rfc7252).

## Getting Started

To get started, either build the `k6` binary or use the
[`goliothiot/xk6-coap`](https://hub.docker.com/repository/docker/goliothiot/xk6-coap)
image.

### Using the OCI Image

The quickest way to get started is by using the `goliothiot/xk6-coap` image,
which is [built](./images/xk6-coap/Dockerfile) from this repository and
published to DockerHub. Tests can be supplied when creating a container from the
image via a [bind mount](https://docs.docker.com/storage/bind-mounts/). For
example, the following command would run the [simple
example](./examples/simple.js) from this repository.

```
docker run -it --rm -v $(pwd)/examples/simple.js:/simple.js goliothiot/xk6-coap k6 run /simple.js --vus 10 --duration 5s
```

### Building a k6 Binary

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
