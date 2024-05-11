package handlers

import (
	"RockPaperScissor/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebsocketGame(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		gameId := c.Param("gameId")
		if gameId == "" {
			c.String(http.StatusBadRequest, "Please provide a Game ID")
			return
		}

		if _, ok := hub.Pools[gameId]; !ok {
			pool := pkg.NewPool()

			hub.Pools[gameId] = pool

			go pool.Start()

		}

		var pool *pkg.Pool = hub.Pools[gameId]

		if len(pool.Clients) == 2 {
			c.String(http.StatusForbidden, "Lobby is Full")
			return
		}

		pkg.ServeWs(pool, c.Writer, c.Request)
	}
}

func HandleNewGame(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameId := pkg.GenerateRandomString()

		c.String(http.StatusOK, gameId)
		return
	}
}
