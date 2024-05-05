package pkg

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Pool struct {
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	board      map[int]*Hand
	gameStatus chan int
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		board:      make(map[int]*Hand),
		gameStatus: make(chan int),
	}
}

type Hand struct {
	client *Client
	hand   string
}

type Hub struct {
	Pools map[string]*Pool
}

func NewHub() *Hub {
	return &Hub{
		Pools: make(map[string]*Pool),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.register:
			pool.Clients[client] = true
			pool.board[client.ID] = &Hand{
				client: client,
				hand:   "X",
			}
			for client, _ := range pool.Clients {
				client.Conn.WriteMessage(1, []byte("New User Joined"))
			}
			break
		case message := <-pool.broadcast:

			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			break

		case client := <-pool.unregister:
			delete(pool.Clients, client)
			delete(pool.board, client.ID)
			break

		case gameStatus := <-pool.gameStatus:
			fmt.Println("in GameStatus", len(pool.board), gameStatus)
			var isReady bool = IsPlayersReady(pool)
			if isReady {
				fmt.Println("Who WIn Logic")
				var player1Hand string = pool.board[1].hand
				var player2Hand string = pool.board[2].hand

				var winnerId int = PlayGame(player1Hand, player2Hand)

				fmt.Println("Hands", player1Hand, player2Hand)

				if winnerId == 0 {
					pool.broadcast <- Message{
						MessageType: 1,
						Message:     "It's a tie",
					}
				} else {
					for c, _ := range pool.Clients {
						if winnerId == c.ID {
							c.Conn.WriteMessage(1, []byte("You Win !"))
							continue
						}
						c.Conn.WriteMessage(1, []byte("You Lose !"))
					}
				}
				pool.board[1].hand = "X"
				pool.board[2].hand = "X"
			} else {
				pool.board[gameStatus].client.Conn.WriteMessage(1, []byte("Waiting for Player 2 ..."))
			}
			// }
			break
		}
	}
}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Conn: conn,
		pool: pool,
		ID:   len(pool.Clients) + 1,
	}

	pool.register <- client

	go client.Read()
}
