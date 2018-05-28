package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan Event

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan Client

	channels map[string]*Channel
}

func (broker *broker) GetChannelSize() int {
	return len(broker.channels)
}

func (broker *broker) GetClientSize() int {
	var clients int
	for _, c := range broker.channels {
		clients += len(c.Clients)
	}
	return clients
}

func (broker *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, ErrorStreamingUnsupported.Error(), http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	var password string
	var channelId string
	var clientId string

	if _, ok := query["channelId"]; !ok || len(query["channelId"][0]) != 12 {
		http.Error(w, ErrorInvalidChannelId.Error(), http.StatusBadRequest)
		return
	} else {
		channelId = query["channelId"][0]
	}

	if _, ok := query["clientId"]; !ok || len(query["clientId"][0]) < 3 {
		http.Error(w, ErrorInvalidClientId.Error(), http.StatusBadRequest)
		return
	} else {
		clientId = query["clientId"][0]
	}

	if _, ok := query["password"]; ok {
		password = query["password"][0]
	}

	if c, ok := broker.channels[channelId]; ok {

		if len(c.Password) != 0 && password != c.Password {

			http.Error(w, ErrorInvalidChannelPassword.Error(), http.StatusUnauthorized)
			return
		}

		for _, client := range c.Clients {
			if strings.EqualFold(client.ClientId, clientId) {
				http.Error(w, ErrorClientIdAlreadyUsed.Error(), http.StatusBadRequest)
				return
			}
		}

	} else {
		broker.channels[channelId] = &Channel{
			channelId,
			password,
			[]Client{},
		}
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	client := Client{
		clientId,
		channelId,
		make(chan Event),
	}

	// Signal the broker that we have a new connection
	broker.newClients <- client

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- client
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.closingClients <- client
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		event := <-client.Channel
		msg, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)

		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()
	}

}

func (broker *broker) listen() {

	for {
		select {
		case client := <-broker.newClients:

			// A new client has connected.
			// Register their message channel

			if c, ok := broker.channels[client.ChannelId]; ok {
				broker.Notifier <- Event{
					client.ChannelId,
					Connect,
					client.ClientId,
					nil,
				}
				c.Clients = append(c.Clients, client)

				var clients []string

				for _, clientInChannel := range c.Clients {
					clients = append(clients, clientInChannel.ClientId)
				}
				client.Channel <- Event{
					client.ChannelId,
					ClientList,
					"server",
					clients,
				}
			} else {
				log.Println(ErrorChannelDoesNotExist)
			}

		case client := <-broker.closingClients:

			if c, ok := broker.channels[client.ChannelId]; ok {

				for i, cl := range c.Clients {
					if cl == client {
						c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
					}
				}

				if len(c.Clients) == 0 {
					delete(broker.channels, client.ChannelId)
				} else {
					broker.channels[client.ChannelId] = c
					broker.Notifier <- Event{
						client.ChannelId,
						Disconnect,
						client.ClientId,
						nil,
					}
				}
			}

		case event := <-broker.Notifier:

			if c, ok := broker.channels[event.ChannelId]; ok {
				for _, client := range c.Clients {
					client.Channel <- event
				}
			}
		}
	}
}
