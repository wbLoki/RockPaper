package api

import (
	"RockPaperScissor/handlers/game"
	"RockPaperScissor/pkg"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	addr string
	pool *pkg.Pool
	rdb  *redis.Client
}

func NewApiServer(addr string, pool *pkg.Pool, rdb *redis.Client) *APIServer {
	return &APIServer{
		addr: addr,
		pool: pool,
		rdb:  rdb,
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

	gameHandler := game.NewHandler(s.pool, s.rdb)
	gameHandler.RegisterRoutes(router)

	log.Println("Listening on ", s.addr)
	return router.Run(s.addr)
}
