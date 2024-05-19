package game

import (
	"RockPaperScissor/config"
	"RockPaperScissor/pkg"
	"RockPaperScissor/types"
	"RockPaperScissor/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	router.PUT("/players", h.HandleUpdatePlayer)
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

	pkg.ServeWs(h.pool, h.rdb, redisGame, c)
}

func (h *Handler) HandleNewGame(c *gin.Context) {

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

func (h *Handler) HandleUpdatePlayer(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid Request")
		return
	}

	var playerPaylaod types.UpdatePlayerPayload
	if err := json.Unmarshal(body, &playerPaylaod); err != nil {
		c.String(http.StatusBadRequest, "Invalid Request")
		return
	}

	var playerRedis types.PlayerRedis
	if err := utils.GetFromRedis(h.rdb, playerPaylaod.Id, &playerRedis); err != nil {
		c.String(http.StatusBadRequest, "Invalid Request")
		return
	}

	playerRedis.Name = playerPaylaod.Name

	if err := utils.SetRedis(h.rdb, playerPaylaod.Id, playerRedis); err != nil {
		c.String(http.StatusInternalServerError, "Ouch")
		return
	}

	c.String(http.StatusOK, "Updated Succ")
}
