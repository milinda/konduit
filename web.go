package main

import (
	"embed"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//go:embed client/build
var embeddedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "client/build")
	if err != nil {
		zap.S().Panic(err)
	}

	return http.FS(fsys)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Panic(err)
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	go client.update()
}

func random(hub *Hub) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Generating random value")
				hub.broadcast <- []byte(string(rune(rand.Intn(100))))
				log.Println("Sent random value")
			}
		}
	}()
}

func startWebApp(hub *Hub) {

	http.Handle("/", http.FileServer(getFileSystem()))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		zap.S().Panic(err)
	}
}
