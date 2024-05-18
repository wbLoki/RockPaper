package pkg

import (
	"RockPaperScissor/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(pool *Pool, rdb *redis.Client, redisGame types.RedisGame, c *gin.Context) {
	w := c.Writer
	r := c.Request
	gameId := c.Param("gameId")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var ctx = context.Background()
	clientId := uuid.New().String()

	redisGame.Lobby = append(redisGame.Lobby, clientId)
	redisGame.Hands[clientId] = "X"
	redisGameMarshaled, redisErr := json.Marshal(redisGame)
	if redisErr != nil {
		fmt.Println(redisErr)
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	if err := rdb.Set(ctx, gameId, string(redisGameMarshaled), 0).Err(); err != nil {
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	var playerRedis types.PlayerRedis = types.PlayerRedis{
		Score:  0,
		GameId: gameId,
	}

	playerRedisMarsheld, _ := json.Marshal(playerRedis)

	if err := rdb.Set(ctx, clientId, string(playerRedisMarsheld), 0).Err(); err != nil {
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	client := &Client{
		Conn:   conn,
		pool:   pool,
		ID:     clientId,
		GameId: gameId,
		gameBoard: &GameBoard{
			score: 0,
		},
	}

	pool.register <- client

	go client.Read()
}
