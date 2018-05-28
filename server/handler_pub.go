package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *server) HandlePost() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var event Event
		err := decoder.Decode(&event)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(event.ChannelId) != 12 {
			http.Error(w, ErrorInvalidChannelId.Error(), http.StatusBadRequest)
			return
		}

		if len(event.Type) == 0 {
			http.Error(w, ErrorInvalidEventType.Error(), http.StatusBadRequest)
			return
		}

		if len(event.OriginId) < 3 {
			http.Error(w, ErrorInvalidOriginId.Error(), http.StatusBadRequest)
			return
		}

		if c, ok := s.Broker.channels[event.ChannelId]; ok {

			if len(c.Password) != 0 {

				authorization := strings.Split(r.Header.Get("authorization"), " ")
				if len(authorization) != 2 || authorization[0] != "Password" || authorization[1] != c.Password {
					http.Error(w, ErrorInvalidChannelPassword.Error(), http.StatusUnauthorized)
					return
				}
			}

			s.Broker.Notifier <- event

		} else {
			http.Error(w, ErrorChannelDoesNotExist.Error(), http.StatusNotFound)
		}
	}
}
