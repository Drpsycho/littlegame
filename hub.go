package main

import (
	"encoding/json"
	"log"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var HubHandler = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true

			for m := range h.connections {
				if m.userinfo.User != "" {
					log.Println("User ", m.userinfo.User, " ", m.userinfo.U_x, " ", m.userinfo.U_y)
					res, _ := json.Marshal(m.userinfo)
					c.send <- res
				}

			}
			log.Println("")

		case c := <-h.unregister:
			c.userinfo.Status = false
			byebye, _ := json.Marshal(c.userinfo)
			log.Println("User ", c.userinfo.User, " removed !")

			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}

			for m := range h.connections {
				m.send <- byebye
			}

		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
