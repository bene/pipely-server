package main

import (
	"github.com/bene/flowengine-api-sdk/middleware"
	"github.com/bene/wetube-server/server"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	server := server.NewServer()

	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	router.HandleFunc("/", server.HandlePost()).Methods("POST")
	router.Handle("/", server.Broker).Methods("GET")

	http.ListenAndServe(":8080", router)
}
