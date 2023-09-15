package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Game struct {
	id        string
	playerOne *bool
	playerTwo *bool
	status    *bool
	board     [6]int
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	connections []*websocket.Conn
	Games       = map[string]Game{}
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func boolPointer(s bool) *bool {
	return &s
}

func assignPlayers(gameToken string) bool {
	val, ok := Games[gameToken]

	if !ok {
		Games[gameToken] = Game{
			id:        gameToken,
			playerOne: boolPointer(true),
		}
		return false
	} else {
		val.playerTwo = boolPointer(true)
		val.status = boolPointer(true)
		Games[gameToken] = val
		return true
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	connections = append(connections, ws)

	var gameToken string = r.URL.Query().Get("token")
	if assignPlayers(gameToken) {
		castMessage([]byte("gameon"), 1)
	}
	check(err)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		check(err)
		fmt.Println(string(p)) // TODO remove this incoming msg
		castMessage(p, messageType)
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func castMessage(message []byte, mType int) {
	for _, connection := range connections {
		connection.WriteMessage(mType, message)
	}
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
