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

func (c *Client) Read() {
	defer func() {
		c.pool.unregister <- c
		c.pool.broadcast <- Message{
			MessageType: GM,
			Message:     "Player Left",
		}
		c.Conn.Close()
	}()
	for {
		var incomingMessage Message
		err := c.Conn.ReadJSON(&incomingMessage)
		if err != nil {
			fmt.Println(err)
			break
		}

		switch incomingMessage.MessageType {
		case Game:
			c.pool.board[c.ID].hand = incomingMessage.Message
			c.pool.gameStatus <- c.ID
			break
		case Chat:
			c.pool.broadcast <- incomingMessage
		}
	}
}
