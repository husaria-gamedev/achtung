package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type GameState string

const (
	Starting GameState = "starting"
	Running  GameState = "running"
	Finished GameState = "finished"
)

type Direction string

const FPS int = 30
const Edge int = 1000

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
	Alive      bool
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

type GameStateMessage struct {
	Players []PlayerStateMessage `json:"p"`
	State   GameState            `json:"s"`
}

var state GameState = "starting"
var nextStateMilis int = 1000 * 5

var upgrader = websocket.Upgrader{} // use default options

var connections = make(map[string]*websocket.Conn)
var gameState = make(map[string]*PlayerState)

func addConnection(id string, conn *websocket.Conn) {
	connections[id] = conn
	// later initialize players on game init
	gameState[id] = &PlayerState{
		PosX:       0,
		PosY:       0,
		Angle:      0,
		Speed:      3,
		Alive:      false,
		AngleSpeed: 0.07,
		Direction:  Forward,
		Id:         len(connections),
	}
	if state == "starting" || state == "finished" {
		initPlayer(gameState[id])
	}
	fmt.Println("New connection", id, len(connections))
}

func initPlayer(p *PlayerState) {
	p.PosX = float64(rand.IntN(Edge-Edge/10) + int(Edge/20))
	p.PosY = float64(rand.IntN(Edge-Edge/10) + int(Edge/20))
	p.Angle = rand.Float64() * math.Pi * 2
	p.Alive = true
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
			tick()
		}
	}
}

func updatePlayerState(p *PlayerState) {
	if !p.Alive {
		return
	}

	deltaX := p.Speed * math.Cos(p.Angle)
	deltaY := p.Speed * math.Sin(p.Angle)

	p.PosX += deltaX
	p.PosY += deltaY

	if p.Direction == Left {
		p.Angle -= p.AngleSpeed
	} else if p.Direction == Right {
		p.Angle += p.AngleSpeed
	}

	if p.PosX < 0 || p.PosX > float64(Edge) || p.PosY < 0 || p.PosY > float64(Edge) {
		p.Alive = false
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
	msg := GameStateMessage{
		make([]PlayerStateMessage, 0, len(connections)),
		state,
	}

	for _, p := range gameState {
		pmsg := getPlayerStateMessage(p)
		msg.Players = append(msg.Players, pmsg)
	}
	return msg
}

func runningStateTick() {
	for _, p := range gameState {
		updatePlayerState(p)
	}

	alive := 0
	for _, p := range gameState {
		if p.Alive {
			alive += 1
		}
	}

	if alive <= 1 {
		state = "finished"
		nextStateMilis = 5 * 1000
		for _, p := range gameState {
			initPlayer(p)
		}
	}
}

func startingStateTick() {
	if len(connections) > 1 {
		nextStateMilis -= int(time.Second / time.Duration(FPS) / 10e5)
	} else {
		nextStateMilis = 3 * 1000
	}
	if nextStateMilis < 0 {
		state = "running"
	}
}

func finishedStateTick() {
	nextStateMilis -= int(time.Second / time.Duration(FPS) / 10e5)
	if nextStateMilis < 0 {
		state = "starting"
		nextStateMilis = 1000 * 5
	}
}

func tick() {
	switch state {
	case "running":
		runningStateTick()
	case "starting":
		startingStateTick()
	case "finished":
		finishedStateTick()
	}

	msg := getGameStateMessage()

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
