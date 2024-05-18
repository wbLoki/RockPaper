package game

import (
	"RockPaperScissor/pkg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TestGame(t *testing.T) {
	rdb := redis.Client{}
	pool := pkg.NewPool(&rdb)
	handler := NewHandler(pool, &rdb)

	t.Run("Should fail if not ok", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/game", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := gin.New()
		router.POST("/game", handler.HandleNewGame)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusMovedPermanently {
			t.Errorf("expected %d got %d", http.StatusMovedPermanently, rr.Code)
		}
	})

	t.Run("Should fail if gameId found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/game/testId/valid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := gin.New()
		router.GET("/game/:gameId/valid", handler.HandleValidGame)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected %d got %d", http.StatusNotFound, rr.Code)
		}
	})
}
