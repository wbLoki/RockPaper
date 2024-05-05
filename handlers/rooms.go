package handlers

import (
	"ChatAppGin/pkg"
	rooms "ChatAppGin/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Room struct {
	ID string `json:"id"`
}

func CreatingRoom(c *gin.Context) {
	var roomId string = rooms.GenerateRandomString()
	var room Room = Room{
		ID: roomId,
	}
	c.JSON(http.StatusOK, room)

}

func HandleWebsocketRoom(hub *pkg.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		roomId := c.Param("roomId")
		if roomId == "" {
			c.String(http.StatusBadRequest, "Please give a room id")
			return
		}

		if _, ok := hub.Pools[roomId]; !ok {
			// fmt.Println("Pool Doesn't Exist Creating New Pool")

			pool := pkg.NewPool()

			hub.Pools[roomId] = pool

			go pool.Start()

		}

		var pool *pkg.Pool = hub.Pools[roomId]

		if len(pool.Clients) == 2 {
			c.String(http.StatusForbidden, "Room is Full")
			return
		}

		pkg.ServeWs(pool, c.Writer, c.Request)
	}
}
