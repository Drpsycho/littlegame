package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type inputmsgData struct {
	Dx int `json:"Dx"`
	Dy int `json:"Dy"`
	Id int `json:"Id"`
}

type inputmsg struct {
	Name   string       `json:"Name"`
	Action string       `json:"Action"`
	Data   inputmsgData `json:"Data"`
	Status bool
}

type actionmsg struct {
	In inputmsg
	Dx *int
	Dy *int
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	userinfo inputmsg
}

func serverWs(rw http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println("websocket working")
	defer ws.Close()

	c := &connection{ws: ws, send: make(chan []byte, 256)}
	HubHandler.register <- c
	go c.writePump()
	c.readPump()
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		HubHandler.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		var dat inputmsg

		if err := json.Unmarshal(message, &dat); err != nil {
			log.Println(err)
			continue
		}
		if dat.Action == "myname" {
			c.userinfo.Status = true
			c.userinfo.Name = dat.Name
			c.userinfo.Data.Dx = dat.Data.Dx
			c.userinfo.Data.Dy = dat.Data.Dy

			res, _ := json.Marshal(inputmsg{
				Name:   c.userinfo.Name,
				Action: "addplayer",
				Data:   inputmsgData{Dx: c.userinfo.Data.Dx, Dy: c.userinfo.Data.Dy},
			})

			log.Println("New user - ", c.userinfo.Name)

			for m := range HubHandler.connections {
				if m.userinfo.Name != c.userinfo.Name {
					m.send <- res
				}
			}

			for m := range HubHandler.connections {
				if m.userinfo.Name != c.userinfo.Name {
					res, _ := json.Marshal(inputmsg{
						Name:   m.userinfo.Name,
						Action: "addplayer",
						Data:   inputmsgData{Dx: m.userinfo.Data.Dx, Dy: m.userinfo.Data.Dy},
					})
					log.Println("User ", m.userinfo.Name, " send to ", c.userinfo.Name)
					c.send <- res
				}
			}
		}

		if dat.Action != "" {
			Action <- actionmsg{
				In: dat,
				Dx: &c.userinfo.Data.Dx,
				Dy: &c.userinfo.Data.Dy,
			}
		}
	}
}

var Action = make(chan actionmsg, 100)

func ActionHandler() {
	for {
		_a := <-Action

		if _a.In.Action == "move" {
			collision(&_a)

			res, _ := json.Marshal(inputmsg{
				Name:   _a.In.Name,
				Action: "moveplayer",
				Data:   inputmsgData{Dx: _a.In.Data.Dx, Dy: _a.In.Data.Dy},
			})

			HubHandler.broadcast <- res
			//send to move handler
		}
		if _a.In.Action == "setbomb" {
			//send to bombs handler
		}
		if _a.In.Action == "setname" {
			//set name
		}

	}
}
