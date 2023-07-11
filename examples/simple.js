import { Client } from 'k6/x/coap';

export default function() {
	// Create new client and connect.
	let client;
	try {
		client = new Client("coap.golioth.io:5684");
	} catch (e) {
		console.log(e);
	}

	// Verify connection.
	try {
		let res = client.get("/hello", 10);
		console.log(String.fromCharCode(...res.body));
	} catch (e) {
		console.log(e);
	}

	// Send data.
	try {
		let res = client.post("/.s", "application/json", '{"hello": "world"}', 10);
		console.log(res.code);
	} catch (e) {
		console.log(e);
	}

	// Get JSON data.
	try {
		let res = client.get("/.u/desired", 10);
		let json = JSON.parse(String.fromCharCode(...res.body));
		console.log(json.sequenceNumber);
	} catch (e) {
		console.log(e);
	}

	// Close connection.
	try {
		client.close();
	} catch (e) {
		console.log(e);
	}
}
