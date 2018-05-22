package server

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
)

type broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan Event

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan Client

	channels map[string][]Client
}

func (broker *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	if _, ok := query["channelId"]; !ok {
		http.Error(w, "No channel specified!", http.StatusBadRequest)
		return
	} else {
		if len(query["channelId"][0]) != 12 {
			http.Error(w, "Invalid channel id!", http.StatusBadRequest)
			return
		}
	}

	if _, ok := query["clientId"]; !ok {
		http.Error(w, "No client id specified!", http.StatusBadRequest)
		return
	} else {
		if len(query["clientId"][0]) < 3 {
			http.Error(w, "Invalid client id! Id must not be shorter than three letters.", http.StatusBadRequest)
			return
		}
	}

	channelId := query["channelId"][0]
	clientId := query["clientId"][0]

	if c, ok := broker.channels[channelId]; ok {

		for _, client := range c {
			if strings.EqualFold(client.ClientId, clientId) {
				http.Error(w, "Invalid client id! Id already in use.", http.StatusBadRequest)
				return
			}
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

		// Flush the data immediatly instead of buffering it for later.
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
					"CONNECT",
					client.ClientId,
					nil,
				}
				broker.channels[client.ChannelId] = append(c, client)

				var clients []string

				for _, clientInChannel := range c {
					clients = append(clients, clientInChannel.ClientId)
				}
				client.Channel <- Event{
					client.ChannelId,
					"CLIENT_LIST",
					"server",
					clients,
				}
			} else {
				broker.channels[client.ChannelId] = []Client{client}
				client.Channel <- Event{
					client.ChannelId,
					"CLIENT_LIST",
					"server",
					[]string{client.ClientId},
				}
			}

			log.Printf("Client added. %d channels", len(broker.channels))
		case client := <-broker.closingClients:

			if c, ok := broker.channels[client.ChannelId]; ok {

				for i, cl := range c {
					if cl == client {
						c = append(c[:i], c[i+1:]...)
					}
				}

				if len(c) == 0 {
					delete(broker.channels, client.ChannelId)
				} else {
					broker.channels[client.ChannelId] = c
					broker.Notifier <- Event{
						client.ChannelId,
						"DISCONNECT",
						client.ClientId,
						nil,
					}
				}
			}

			log.Printf("Client removed. %d channels", len(broker.channels))
		case event := <-broker.Notifier:

			if c, ok := broker.channels[event.ChannelId]; ok {
				for _, client := range c {
					client.Channel <- event
				}
			}
		}
	}
}
