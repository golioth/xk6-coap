# xk6-coap

`xk6-coap` is a [`k6` extension](https://k6.io/docs/extensions/) for the
[Constrained Application Protocol
(CoAP)](https://www.rfc-editor.org/rfc/rfc7252).

> **MATURITY NOTICE**: `xk6-coap` is under active development and breaking API
> changes are to be expected. The project has been open sourced early to ensure
> that external consumers have an opportunity to influence and contribute to its
> future direction.

## Getting Started

To get started, either build the `k6` binary or use the
[`golioth/xk6-coap`](https://hub.docker.com/repository/docker/golioth/xk6-coap)
image.

### Using the OCI Image

The quickest way to get started is by using the `golioth/xk6-coap` image,
which is [built](./images/xk6-coap/Dockerfile) from this repository and
published to DockerHub. Tests can be supplied when creating a container from the
image via a [bind mount](https://docs.docker.com/storage/bind-mounts/). For
example, the following command would run the [simple
example](./examples/simple.js) from this repository.

```
docker run -it --rm -e COAP_PSK_ID=<YOUR-PSK-ID> -e COAP_PSK=<YOUR-PSK> -v $(pwd)/examples/simple.js:/simple.js golioth/xk6-coap k6 run /simple.js --vus 10 --duration 5s
```

`xk6-coap` supports authentication via pre-shared keys (PSKs) and client
certificates. The former is provided by specifying environment variables, while
the latter is provided by specifying file paths. The `simple.js` example passes
`COAP_PSK_ID` and `COAP_PSK` to the instantiated `Client`, which will cause it
use use the respective values for PSK authentication. If a `Client` is
instantiated with both PSK environment variables and certificate file paths,
certificate authentication will take precedence.

```js
client = new Client(
	"coap.golioth.io:5684",
	"COAP_PSK_ID",
	"COAP_PSK",
);
```
> Only PSK is provided and it will be used for authentication.

```js
client = new Client(
	"coap.golioth.io:5684",
	"",
	"",
	"path/to/client/crt.pem",
	"path/to/client/key.pem",
);
```
> Only certificate is provided and it will be used for authentication.

```js
client = new Client(
	"coap.golioth.io:5684",
	"COAP_PSK_ID",
	"COAP_PSK",
	"path/to/client/crt.pem",
	"path/to/client/key.pem",
);
```
> Both are provided but certificate takes precedence.


### Building a k6 Binary

Using `k6` extentions requires including the extension(s) in a `k6` build. The
[`xk6`](https://github.com/grafana/xk6) tool will handle executing the build,
and the `Makefile` in this repository will ensure that `xk6` is installed and
produces a valid build.

```
make build
```

This will produce a `k6` binary in the current working directory. To execute a
test with 2 virtual users making connections and sending `GET` messages for 5
seconds, the following command could be run.

```
./k6 run ./examples/simple.js --vus 10 --duration 5s
```

Reference the [`k6` documentation](https://k6.io/docs/using-k6/test-lifecycle/)
for more information on how to configure and run tests.

## Attribution

`xk6-coap` is essentially glue machinery that allows for
[`plgd-dev/go-coap`](https://github.com/plgd-dev/go-coap) /
[`pion/dtls`](https://github.com/plgd-dev/go-coap) functionality to be exposed
to `k6` tests. This project would not be possible without the work done by
contributors (some of whom are on the Golioth team!) on both of those projects.
