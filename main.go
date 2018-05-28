package main

import (
	"github.com/bene/flowengine-api-sdk/middleware"
	"github.com/bene/pipely-server/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

const banner = `

      ___                       ___           ___           ___       ___     
     /\  \          ___        /\  \         /\  \         /\__\     |\__\    
    /::\  \        /\  \      /::\  \       /::\  \       /:/  /     |:|  |   
   /:/\:\  \       \:\  \    /:/\:\  \     /:/\:\  \     /:/  /      |:|  |   
  /::\~\:\  \      /::\__\  /::\~\:\  \   /::\~\:\  \   /:/  /       |:|__|__ 
 /:/\:\ \:\__\  __/:/\/__/ /:/\:\ \:\__\ /:/\:\ \:\__\ /:/__/        /::::\__\
 \/__\:\/:/  / /\/:/  /    \/__\:\/:/  / \:\~\:\ \/__/ \:\  \       /:/~~/~   
      \::/  /  \::/__/          \::/  /   \:\ \:\__\    \:\  \     /:/  /     
       \/__/    \:\__\           \/__/     \:\ \/__/     \:\  \    \/__/      
                 \/__/                      \:\__\        \:\__\              
                                             \/__/         \/__/
`

func main() {

	log.Println(banner)

	address := ":6550"
	if addr, ok := os.LookupEnv("ADDRESS"); ok && len(addr) != 0 {
		address = addr
	}
	log.Printf("Using server address: %s", address)

	server := server.NewServer()

	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	router.HandleFunc("/publish", server.CreateHandlerPublish()).Methods("POST")
	router.HandleFunc("/channels", server.CreateHandlerChannels()).Methods("GET")
	router.HandleFunc("/channel/{channelId}", server.CreateHandlerChannel()).Methods("GET")
	router.Handle("/subscribe", server.Broker).Methods("GET")

	go func() {
		for {
			log.Printf("Open channels: %d", server.Broker.GetChannelSize())
			log.Printf("Connected clients: %d", server.Broker.GetClientSize())
			time.Sleep(time.Minute * 5)
		}
	}()

	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
