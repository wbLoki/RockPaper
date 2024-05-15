package tests

import (
	"RockPaperScissor/handlers"
	"RockPaperScissor/pkg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	hub := pkg.NewHub()
	hub.Pools["testOk"] = pkg.NewPool()
	hub.Pools["testFull"] = pkg.NewPool()

	hub.Pools["testFull"].Clients[&pkg.Client{
		ID: 1,
	}] = true
	hub.Pools["testFull"].Clients[&pkg.Client{
		ID: 2,
	}] = true

	r := gin.Default()
	r.GET("/game/:gameId", handlers.HandleWebsocketGame(hub))
	r.POST("/game", handlers.HandleNewGame(hub))
	r.GET("/game/:gameId/valid", handlers.HandleValidGame(hub))
	return r
}

func TestWebSocketHandler(t *testing.T) {

	r := setupRouter()

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + srv.URL[4:] + "/game/testOk"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not connect to WebSocket: %v", err)
	}
	defer ws.Close()

	assert.NotNil(t, ws)

}

func TestNewGameHandler(t *testing.T) {
	r := setupRouter()
	req, err := http.NewRequest("POST", "/game", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestValidGameNotFound(t *testing.T) {
	r := setupRouter()
	gameId := "testId"
	req, err := http.NewRequest("GET", "/game/"+gameId+"/valid", nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}
func TestValidGameFound(t *testing.T) {
	r := setupRouter()
	gameId := "testOk"
	req, err := http.NewRequest("GET", "/game/"+gameId+"/valid", nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
func TestValidGameFull(t *testing.T) {
	r := setupRouter()
	gameId := "testFull"
	req, err := http.NewRequest("GET", "/game/"+gameId+"/valid", nil)

	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

}
