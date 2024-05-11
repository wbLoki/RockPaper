package handlers

import (
	"RockPaperScissor/pkg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebsocketGame(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		gameId := c.Param("gameId")

		if _, ok := hub.Pools[gameId]; !ok || gameId == "" {
			c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
			return
		}

		var pool *pkg.Pool = hub.Pools[gameId]

		if len(pool.Clients) == 2 {
			c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
			return
		}

		pkg.ServeWs(pool, c.Writer, c.Request)
	}
}

func HandleNewGame(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameId := pkg.GenerateRandomString()

		pool := pkg.NewPool()
		hub.Pools[gameId] = pool

		go pool.Start()

		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://localhost:3000/game/%s", gameId))

	}
}

func HandleValidGame(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameId := c.Param("gameId")

		if _, ok := hub.Pools[gameId]; !ok {
			c.String(http.StatusNotFound, "Game not found !")
			return
		}

		if pool := hub.Pools[gameId]; len(pool.Clients) >= 2 {
			c.String(http.StatusForbidden, "Game is Full")
			return
		}
		c.String(http.StatusOK, "Game Found !")

	}
}
