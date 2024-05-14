package game

import (
	"RockPaperScissor/pkg"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	hub *pkg.Hub
}

func NewHandler(hub *pkg.Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/game", h.HandleNewGame)
	router.GET("/game/:gameId", h.HandleWebsocketGame)
	router.GET("/game/:gameId/valid", h.HandleValidGame())
}

var clientUrl = os.Getenv("CLIENT_URL")

func (h *Handler) HandleWebsocketGame(c *gin.Context) {

	gameId := c.Param("gameId")

	if _, ok := h.hub.Pools[gameId]; !ok || gameId == "" {
		c.Redirect(http.StatusPermanentRedirect, clientUrl)
		return
	}

	var pool *pkg.Pool = h.hub.Pools[gameId]

	if len(pool.Clients) == 2 {
		c.Redirect(http.StatusPermanentRedirect, clientUrl)
		return
	}

	pkg.ServeWs(pool, c.Writer, c.Request)
}

func (h *Handler) HandleNewGame(c *gin.Context) {
	gameId := pkg.GenerateRandomString()
	fmt.Println("New Game")
	pool := pkg.NewPool()
	h.hub.Pools[gameId] = pool

	go pool.Start()

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%sgame/%s", clientUrl, gameId))

}

func (h *Handler) HandleValidGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		gameId := c.Param("gameId")

		if _, ok := h.hub.Pools[gameId]; !ok {
			c.String(http.StatusNotFound, "Game not found !")
			return
		}

		if pool := h.hub.Pools[gameId]; len(pool.Clients) >= 2 {
			c.String(http.StatusForbidden, "Game is Full")
			return
		}
		c.String(http.StatusOK, "Game Found !")

	}
}
