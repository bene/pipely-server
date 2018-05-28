package server

import "log"

var builtin_events string

func r() {
	builtin_events = ""
	log.Println(builtin_events)
}
