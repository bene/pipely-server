# Pipely

Build connected apps without the need to write a server. Pipely enables you to connect clients and exchange data with a server-client/client-client hybrid architecture.
The Pipely server is a lightweight, concurrent SSE server with allows clients to created temporary channels which can be joined by other clients. The channels can be created with or without a password protection.

### Getting started

In order to begin writing serverless connected apps you need to:
- Deploy a Pipely server
- Create/Join a channel
- Exchange data

## Deploy a Pipely server

In order to setup a Pipely server run:
```sh
$ git clone git@github.com:bene/pipely-server.git
$ cd pipely-server
$ docker build -t pipely/server .
```
After a successful build start a Docker container:
```sh
$ docker run --name="pipely-1" -p 6550:6550 --restart=always pipely/server
```

## Create/Join a channel

Once a server is up and running you can subscribe with an EventSource or SSE client. This example uses the default JavaScript EventSource client, with is currently supported by Chrome, Firefox, Safari and Opera. Polyfills are available for Edge and IE.
```javascript
const eventSource = new EventSource('//localhost:6550/subscribe?channelId=CHANNEL_ID&clientId=CLIENT_ID&password=CHANNEL_PASSWORD')
```
| Parameter        | Type        | Criteria                                                                                             | Required |
| ---------------- |:-----------:| ---------------------------------------------------------------------------------------------------- |:--------:|
| channelId        | String      | Twelve characters long                                                                               | Yes      |
| password         | String      | Only needed if the channel has or should have a password                                             | No       |
| clientId         | String      | Must be at least three characters long and has to be unique in channel                               | Yes      |

If a client connects to a non existent channel, the channel will be created, with or without a password, depending on the query. If the channel exists, the password has to be valid. When the last member of a channel disconnects, the channel will be destroyed by the server, and anyone can re-create it with a new password.

## Exchange data
To exchange data, events are used. An event is just a JSON object with specific fields. In order to publish an event to a channel, a POST request has to be sent to the server, in this case cURL is used:
```sh
curl -X POST \
  http://localhost:6550/publish \
  -H 'Authorization: Password CHANNEL_PASSWORD' \
  -H 'Content-Type: application/json' \
  -d '{
	"channel_id":"CHANNEL_ID",
	"type": "EVENT_TYPE",
	"origin_id": "ORIGIN_ID",
	"payload": {
		"test":"Hello World!"
	}
}'
```

#### Request Headers:
| Field            | Type        | Description                                 | Value                     |
| ---------------- |:-----------:| ------------------------------------------- | ------------------------- |
| Authorization    | String      | Needed if the channel has a password        | Password CHANNEL_PASSWORD |
| Content-Type     | String      | Needed since the body is a JSON object this | application/json          |

#### Request Body:
| Field            | Type        | Criteria                                                                                             |
| ---------------- |:-----------:| ---------------------------------------------------------------------------------------------------- |
| channel_id       | String      | Twelve characters long                                                                               |
| type             | String      | Must be at lease on character long                                                                   |
| origin_id        | String      | Must be at least three characters long, should either be 'server' or a client id by a channel member |
| payload          | JSON Object | Must be a valid JSON object or can be undefined                                                      |

To receive data the SSE clients have to manage incoming messages:
```javascript
eventSource.onmessage = (e) => {
    const data = JSON.parse(e.data);

    const eventType = data['type'];
    const originId = data['origin_id'];
    const payload = data['payload'];

    ...
}

```