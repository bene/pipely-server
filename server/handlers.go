package server

import (
	"net/http"
	"encoding/json"
)

func (s *server) HandlePost() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var event Event
		err := decoder.Decode(&event)

		if err != nil {
			http.Error(w, "Invalid event message!", http.StatusBadRequest)
			return
		}

		s.Broker.Notifier <- event
	}
}