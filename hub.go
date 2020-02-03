package main

import "encoding/json"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.clients[conn] = true
			jsonMsg, _ := json.Marshal(Message{conn.id, "A new client has connected."})
			h.send(jsonMsg, conn)
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				close(conn.send)
				delete(h.clients, conn)
				jsonMsg, _ := json.Marshal(Message{conn.id, "A connect has disconnected."})
				h.send(jsonMsg, conn)
			}
		case msg := <-h.broadcast:
			for conn := range h.clients {
				select {
				case conn.send <- msg:
				default:
					close(conn.send)
					delete(h.clients, conn)
				}
			}
		}
	}
}

func (h *Hub) send(msg []byte, ignore *Client) {
	for conn := range h.clients {
		if conn != ignore {
			conn.send <- msg
		}
	}
}
