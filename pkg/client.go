package pkg

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	Conn *websocket.Conn
	pool *Pool
}

// Message Formate {"type":2,"message":"paper"}
type Message struct {
	MessageType int    `json:"type"`
	Message     string `json:"message"`
}

const (
	Chat     = 1
	Game     = 2
	GM       = 3
	Rock     = "rock"
	Paper    = "paper"
	Scissors = "scissors"
)

func (client *Client) Read() {
	defer func() {
		client.pool.unregister <- client
		client.pool.broadcast <- Message{
			MessageType: GM,
			Message:     "Player Left",
		}
		client.Conn.Close()
	}()
	for {
		var incomingMessage Message
		err := client.Conn.ReadJSON(&incomingMessage)
		if err != nil {
			fmt.Println(err)
			break
		}

		switch incomingMessage.MessageType {
		case Game:
			client.pool.board[client.ID].hand = incomingMessage.Message
			client.pool.gameStatus <- client.ID
			break
		case Chat:
			client.pool.broadcast <- incomingMessage
		}
	}
}
