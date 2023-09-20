package main

import (
	function "RockPaper/game"
	"RockPaper/models"
	"RockPaper/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	BluePrint   = []models.Dict{
		{1: "R", 2: "P", 3: "P"},
	}
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
	var gameStatus bool = function.AssignPlayers(gameToken, Games)
	if gameStatus {
		utils.CastMessage([]byte("gameon"), 1, connections)
	}
	check(err)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		check(err)
		// move:1:token
		if strings.Contains(string(p), "move") {
			// move this whole code below into a function
			var GameDetails []string = strings.Split(string(p), ":")
			move, token := GameDetails[1], GameDetails[2]
			playerID, err := strconv.Atoi(string(token[len(token)-1]))
			p_move, err := strconv.Atoi(move)
			check(err)
			token = token[:len(token)-1]
			g := Games[token]
			g.Board[playerID-1] = p_move
			g.Hands[playerID-1] = BluePrint[0][p_move].(string)
			Games[token] = g
			// code below the RuleScore function
		}
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
