package main

import (
	"ChatAppGin/handlers"
	"ChatAppGin/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func main() {
	hub := pkg.NewHub()

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.GET("/", HelloWorldHandler)
	router.GET("/rooms/:roomId", handlers.HandleWebsocketRoom(hub))
	router.Run(":8080")
}
