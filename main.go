package main

import (
	function "RockPaper/game"
	"RockPaper/models"
	"RockPaper/utils"
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

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	connections []*websocket.Conn
	Games       = map[string]models.Game{}
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	connections = append(connections, ws)

	var gameToken string = r.URL.Query().Get("token")
	var gameTokenLength int = len(gameToken)
	var playerID string = string(gameToken[gameTokenLength-1])
	gameToken = gameToken[:gameTokenLength-1]
	fmt.Println(gameToken, playerID)
	if function.AssignPlayers(gameToken, Games) {
		utils.CastMessage([]byte("gameon"), 1, connections)
	}
	check(err)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		check(err)
		fmt.Println(string(p)) // TODO remove this incoming msg
		utils.CastMessage(p, messageType, connections)
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
