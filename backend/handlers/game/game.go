package game

import (
	"RockPaperScissor/config"
	"RockPaperScissor/pkg"
	"RockPaperScissor/types"
	"RockPaperScissor/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	pool *pkg.Pool
	rdb  *redis.Client
}

func NewHandler(pool *pkg.Pool, rdb *redis.Client) *Handler {
	return &Handler{
		pool: pool,
		rdb:  rdb,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/game", h.HandleNewGame)
	router.GET("/game/:gameId", h.HandleWebsocketGame)
	router.GET("/game/:gameId/valid", h.HandleValidGame)
}

func (h *Handler) HandleWebsocketGame(c *gin.Context) {

	gameId := c.Param("gameId")
	var ctx = context.Background()
	var redisGame types.RedisGame

	val, err := h.rdb.Get(ctx, gameId).Result()
	if err != nil {
		c.String(http.StatusNotFound, "Game not found")
		return
	}

	if err := json.Unmarshal([]byte(val), &redisGame); err != nil {
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	if len(redisGame.Lobby) >= 2 {
		c.String(http.StatusForbidden, "Game is Full")
		return
	}

	// if _, ok := h.hub.Pools[gameId]; !ok || gameId == "" {
	// c.Redirect(http.StatusPermanentRedirect, config.Envs.ClientUrl)
	// return
	// }

	// var pool *pkg.Pool = h.hub.Pools[gameId]

	// if len(pool.Clients) == 2 {
	// c.Redirect(http.StatusPermanentRedirect, config.Envs.ClientUrl)
	// return
	// }

	pkg.ServeWs(h.pool, h.rdb, redisGame, c)
}

func (h *Handler) HandleNewGame(c *gin.Context) {
	// Generate new gameId and store it to redis
	// with empty string

	var ctx = context.Background()
	gameId := utils.GenerateRandomString()

	redisGame := types.RedisGame{
		ID:    gameId,
		Lobby: make([]string, 0),
		Hands: make(map[string]string),
	}
	newGame, _ := json.Marshal(redisGame)

	err := h.rdb.Set(ctx, gameId, string(newGame), 0).Err()
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%sgame/%s", config.Envs.ClientUrl, gameId))

}

func (h *Handler) HandleValidGame(c *gin.Context) {
	// gameId will have a key string equal to
	// player1Id, player2Id

	gameId := c.Param("gameId")
	var ctx = context.Background()
	value, err := h.rdb.Get(ctx, gameId).Result()

	if err != nil {
		c.String(http.StatusNotFound, "Game not found !")
		return
	}

	var redisGame types.RedisGame

	if err := json.Unmarshal([]byte(value), &redisGame); err != nil {
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	if len(redisGame.Lobby) >= 2 {
		c.String(http.StatusForbidden, "Game is Full")
		return
	}
	c.String(http.StatusOK, "Game Found !")

}
