package main

import (
	"log"
	"net/http"

	server "github.com/alexmorten/spacy-server"
	"github.com/gorilla/websocket"
)

var gamePool *server.GamePool

func main() {
	gamePool = server.NewGamePool()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	gamePool.AddConnection(conn)
}
