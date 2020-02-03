package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func wsPage(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("cannot generate random uuid")
		return
	}
	client := &Client{id: id.String(), hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.read()
	go client.write()
}
