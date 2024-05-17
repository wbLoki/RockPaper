package pkg

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var clientId int = len(pool.Clients) + 1
	_, ok := pool.board[clientId]
	if ok {
		clientId = 1
	}

	client := &Client{
		Conn: conn,
		pool: pool,
		ID:   clientId,
		gameBoard: &GameBoard{
			score: 0,
		},
	}

	pool.register <- client

	go client.Read()
}