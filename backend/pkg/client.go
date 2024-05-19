package pkg

import (
	"RockPaperScissor/types"
	"RockPaperScissor/utils"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	pool   *Pool
	GameId string
}

const (
	Chat     = 1 // Normal Chat
	Game     = 2 // Game
	GM       = 3 // Game Master
	GE       = 4 // Game End
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
			GameId:      c.GameId,
			ClientId:    c.ID,
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
			var redisGame types.RedisGame

			if err := utils.GetFromRedis(c.pool.rdb, c.GameId, &redisGame); err != nil {
				fmt.Println(err)
			}
			redisGame.Hands[c.ID] = incomingMessage.Message

			if err := utils.SetRedis(c.pool.rdb, c.GameId, redisGame); err != nil {
				fmt.Println(err)
			}

			c.pool.gameId <- c.GameId
		default:
			c.pool.broadcast <- incomingMessage
		}
	}
}
