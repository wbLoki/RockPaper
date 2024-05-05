package pkg

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
			client.Conn.WriteJSON(Message{
				MessageType: GM,
				Message:     "Welcome Player " + strconv.Itoa(client.ID),
			})

			for _client, _ := range pool.Clients {
				_client.Conn.WriteJSON(Message{
					MessageType: GM,
					Message:     "Player " + strconv.Itoa(client.ID) + " Joined",
				})
			}
			pool.Clients[client] = true
			pool.board[client.ID] = &Hand{
				client: client,
				hand:   "X",
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
			var isReady bool = IsPlayersReady(pool)
			if isReady {
				var player1Hand string = pool.board[1].hand
				var player2Hand string = pool.board[2].hand

				var winnerId int = PlayGame(player1Hand, player2Hand)

				if winnerId == 0 {
					for c, _ := range pool.Clients {
						c.Conn.WriteJSON(Message{
							MessageType: GM,
							Message:     "It's a Tie !!",
						})
					}
				} else {
					for c, _ := range pool.Clients {
						if winnerId == c.ID {
							c.Conn.WriteJSON(Message{
								MessageType: GM,
								Message:     "You Win !",
							})
							continue
						}
						c.Conn.WriteJSON(Message{
							MessageType: GM,
							Message:     "You Lose !",
						})
					}
				}
				pool.board[1].hand = "X"
				pool.board[2].hand = "X"
			} else {
				pool.board[gameStatus].client.Conn.WriteJSON(Message{
					MessageType: GM,
					Message:     "Waiting for Player 2 ...",
				})
			}
			break
		}
	}
}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var clientId int = len(pool.Clients) + 1
	_, ok := pool.board[clientId]
	if ok {
		clientId = 1
	}

	client := &Client{
		Conn: conn,
		pool: pool,
		ID:   clientId,
	}

	pool.register <- client

	go client.Read()
}
