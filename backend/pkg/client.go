package pkg

import (
	"RockPaperScissor/types"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        int
	Conn      *websocket.Conn
	pool      *Pool
	gameBoard *GameBoard
}

type GameBoard struct {
	score int
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
		c.pool.broadcast <- types.Message{
			MessageType: GM,
			Message:     "Player Left",
		}
		c.Conn.Close()
	}()
	for {
		var incomingMessage types.Message
		err := c.Conn.ReadJSON(&incomingMessage)
		if err != nil {
			fmt.Println(err)
			break
		}

		switch incomingMessage.MessageType {
		case Game:
			c.pool.board[c.ID].hand = incomingMessage.Message
			c.pool.gameStatus <- c.ID
		default:
			c.pool.broadcast <- incomingMessage
		}
	}
}
