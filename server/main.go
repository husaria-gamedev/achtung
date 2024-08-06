package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameState string

const (
	Waiting GameState = "waiting"
	Running GameState = "running"
)

type Direction string

const (
	Left  Direction = "left"
	Right Direction = "right"
)

var GAME_STATE GameState

var upgrader = websocket.Upgrader{} // use default options

var connections = make(map[string]*websocket.Conn)
var directions = make(map[string]*string)

func addConnection(id string, conn *websocket.Conn) {
	fmt.Println("New connection", id)
	connections[id] = conn
}

func removeConnection(id string) {
	fmt.Println("Deleted connection", id)
	delete(connections, id)
}

func handle_websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	connID := conn.RemoteAddr().String()
	addConnection(connID, conn)
	defer removeConnection(connID)

	// Handle messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// fmt.Printf("Received message from %s: %s\n", connID, string(message))
	}
}

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)
	http.HandleFunc("/websocket", handle_websocket)

	fmt.Println("Listening at port 8000")
	http.ListenAndServe(":8000", nil)
}
