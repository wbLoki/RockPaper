package tests

import (
	"RockPaperScissor/handlers"
	"RockPaperScissor/pkg"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketHandler(t *testing.T) {
	hub := pkg.NewHub()
	r := gin.Default()
	r.GET("/game/:gameId", handlers.HandleWebsocketGame(hub))

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + srv.URL[4:] + "/game/123"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not connect to WebSocket: %v", err)
	}
	defer ws.Close()

	assert.NotNil(t, ws)
}
