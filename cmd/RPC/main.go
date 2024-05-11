package main

import (
	"RockPaperScissor/handlers"
	"RockPaperScissor/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	hub := pkg.NewHub()

	router := gin.New()
	router.Use(CORSMiddleware())
	router.SetTrustedProxies(nil)
	router.GET("/", HelloWorldHandler)
	router.GET("/game/:gameId", handlers.HandleWebsocketGame(hub))
	router.GET("/game/:gameId/valid", handlers.HandleValidGame(hub))
	router.POST("/game", handlers.HandleNewGame(hub))
	router.Run(":8080")
}
