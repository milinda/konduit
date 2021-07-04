package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	webClients    map[*Client]bool
	register      chan *Client
	unregister    chan *Client
	broadcast     chan []byte
	notifications chan map[string]interface{}
}

func newHub() *Hub {
	return &Hub{
		broadcast:     make(chan []byte),
		webClients:    make(map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		notifications: make(chan map[string]interface{}),
	}
}

func (h *Hub) start() {
	for {
		select {
		case client := <-h.register:
			h.webClients[client] = true
		case client := <-h.unregister:
			if _, ok := h.webClients[client]; ok {
				delete(h.webClients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println(fmt.Sprintf("Got a message: %s", message))
			for client := range h.webClients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.webClients, client)
				}
			}

		}
	}
}

func (c *Client) update() {
	ticker := time.NewTicker(80 * time.Second)
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			log.Println(fmt.Sprintf("Writing message: %s", message))

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
