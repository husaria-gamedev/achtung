package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type GameState string

const (
	Waiting GameState = "waiting"
	Running GameState = "running"
)

type Direction string

const FPS int64 = 30

const (
	Left    Direction = "left"
	Right   Direction = "right"
	Forward Direction = "none"
)

type PlayerState struct {
	Id         int
	PosX       float64 // Player's X position
	PosY       float64 // Player's Y position
	Angle      float64
	Speed      float64
	AngleSpeed float64
	Direction  Direction
}

type PlayerStateMessage struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	I int     `json:"i"`
}

type DirectionMessage struct {
	Direction string `json:"d"`
}

type GameStateMessage []PlayerStateMessage

var GAME_STATE GameState

var upgrader = websocket.Upgrader{} // use default options

var connections = make(map[string]*websocket.Conn)
var gameState = make(map[string]*PlayerState)

func addConnection(id string, conn *websocket.Conn) {
	connections[id] = conn
	// later initialize players on game init
	gameState[id] = &PlayerState{
		PosX:       200,
		PosY:       200,
		Angle:      0,
		Speed:      3,
		AngleSpeed: 0.07,
		Direction:  Forward,
		Id:         len(connections),
	}
	fmt.Println("New connection", id, len(connections))
}

func removeConnection(id string) {
	fmt.Println("Deleted connection", id)
	delete(connections, id)
	delete(gameState, id)
}

func parseMessage(msg []byte) *DirectionMessage {
	var dirMessage DirectionMessage
	err := json.Unmarshal(msg, &dirMessage)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return &dirMessage
}

func handleMessage(id string, msg []byte) {
	dirMessage := parseMessage(msg)

	switch dirMessage.Direction {
	case "l":
		gameState[id].Direction = Left
	case "r":
		gameState[id].Direction = Right
	default:
		gameState[id].Direction = Forward
	}
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
		handleMessage(connID, message)
	}
}

func gameLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(FPS))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			update()
		}
	}
}

func updatePlayerState(p *PlayerState) {
	deltaX := p.Speed * math.Cos(p.Angle)
	deltaY := p.Speed * math.Sin(p.Angle)

	p.PosX += deltaX
	p.PosY += deltaY

	if p.Direction == Left {
		p.Angle -= p.AngleSpeed
	} else if p.Direction == Right {
		p.Angle += p.AngleSpeed
	}
}

func getPlayerStateMessage(p *PlayerState) PlayerStateMessage {
	return PlayerStateMessage{
		X: p.PosX,
		Y: p.PosY,
		I: p.Id,
	}
}

func getGameStateMessage() GameStateMessage {
	msg := make(GameStateMessage, 0, len(connections))

	for _, p := range gameState {
		pmsg := getPlayerStateMessage(p)
		msg = append(msg, pmsg)
	}
	return msg
}

func update() {
	for _, p := range gameState {
		updatePlayerState(p)
	}

	msg := getGameStateMessage()
	if len(msg) == 0 {
		return
	}

	bytes, err := json.Marshal(msg)

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, conn := range connections {
		conn.WriteMessage(websocket.TextMessage, bytes)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/websocket", handle_websocket)

	go gameLoop()

	fmt.Println("Listening at port 8000")
	http.ListenAndServe(":8000", nil)
}
