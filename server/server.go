package server

type server struct {
	Broker *broker
}

func NewServer() *server {

	broker := &broker{
		Notifier:       make(chan Event, 1),
		newClients:     make(chan Client),
		closingClients: make(chan Client),
		channels:       make(map[string][]Client),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return &server{
		broker,
	}
}
