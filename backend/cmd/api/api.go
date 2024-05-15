package api

import (
	"RockPaperScissor/handlers/game"
	"RockPaperScissor/pkg"
	"log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	addr string
	hub  *pkg.Hub
}

func NewApiServer(addr string, hub *pkg.Hub) *APIServer {
	return &APIServer{
		addr: addr,
		hub:  hub,
	}
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

func (s *APIServer) Run() error {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.SetTrustedProxies(nil)

	gameHandler := game.NewHandler(s.hub)
	gameHandler.RegisterRoutes(router)

	log.Println("Listening on ", s.addr)
	return router.Run(s.addr)
}
