import { Client } from 'k6/x/coap';

export default function () {
  let client = new Client("coap.golioth.io:5684");
  client.get("/test", 10);
  client.close();
}
