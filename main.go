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
	playerOne *string
	playerTwo *string
	board     [6]int
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	connections []*websocket.Conn
	Games       = map[int]Game{}
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func StringPointer(s string) *string {
	return &s
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)

	// TODO
	/*
		this code below should go into a function in a different endpoint normal http request
		send request to create game and get game id then send it to the ws and wait for player2
		after the player2 joins with the game id , now you can start the game

		id := rand.Intn(math.MaxInt8) // Maxint8 since i don't think im gonna have more than 64 games
		Games[id] = Game{
			playerOne: StringPointer("playerOneTEST"),
			playerTwo: StringPointer("PlayerTwoTEST"),
			id:        "323"}

	*/
	connections = append(connections, ws)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		check(err)
		fmt.Println(string(p))
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
