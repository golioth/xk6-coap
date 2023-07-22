import { fail } from 'k6';
import { setTimeout } from "k6/x/timers"
import { Client } from 'k6/x/coap';


export default function() {
	// Create new client and connect.
	let client;
	try {
		client = new Client(
			"coap.golioth.io:5684",
			"COAP_PSK_ID",
			"COAP_PSK",
			// path/to/client/crt.pem,
			// path/to/client/key.pem,
		);
	} catch (e) {
		fail(e);
	}

	// Verify connection.
	try {
		let res = client.get("/hello", 10);
		console.log(String.fromCharCode(...res.body));
	} catch (e) {
		fail(e);
	}

	// Send data.
	try {
		let res = client.post("/.s", "application/json", '{"hello": "world"}', 10);
		console.log(res.code);
	} catch (e) {
		fail(e);
	}

	// Get JSON data.
	try {
		let res = client.get("/.u/desired", 10);
		let json = JSON.parse(String.fromCharCode(...res.body));
		console.log(json.sequenceNumber);
	} catch (e) {
		fail(e);
	}

	// Start RPC observation.
	try {
		client.observe("/.rpc", 15, (req) => {
			let json;
			try {
				json = JSON.parse(String.fromCharCode(...req.body));
			} catch (e) {
				// First message is acknowledgement of
				// observation.
				console.log(e);
				return;
			}
			try {
				console.log(json);
				client.post("/.rpc/status", "application/json", '{"id": "' + json.id + '", "statusCode": 0, "detail":"ack"}', 10);
			} catch (e) {
				fail(e);
			}
		});
	} catch (e) {
		fail(e);
	}

	// Start OTA observation.
	try {
		client.observe("/.u/desired", 15, (req) => {
			let json;
			try {
				json = JSON.parse(String.fromCharCode(...req.body));
			} catch (e) {
				return;
			}
			console.log(json);
		});
	} catch (e) {
		fail(e);
	}

	// Wait for observations to complete.
	setTimeout(() => {
		// Close connection.
		try {
			client.close();
		} catch (e) {
			fail(e);
		}
	}, 20000)
}
