package pkg

import (
	"fmt"
	"strconv"
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

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.register:
			pool.handleClientRegistration(client)
		case message := <-pool.broadcast:
			pool.broadcastMessage(message)
		case client := <-pool.unregister:
			pool.unregisterClient(client)
		case gameStatus := <-pool.gameStatus:
			pool.handleGameStatus(gameStatus)
		}
	}
}

func (pool *Pool) handleClientRegistration(client *Client) {
	client.Conn.WriteJSON(Message{
		MessageType: GM,
		Message:     "Welcome Player " + strconv.Itoa(client.ID),
	})

	for _client := range pool.Clients {
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
}

func (pool *Pool) broadcastMessage(message Message) {
	for client := range pool.Clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (pool *Pool) unregisterClient(client *Client) {
	delete(pool.Clients, client)
	delete(pool.board, client.ID)
}

func (pool *Pool) handleGameStatus(gameStatus int) {
	isReady := IsPlayersReady(pool)
	if isReady {
		pool.playGame()
	} else {
		pool.notifyWaitingPlayers(gameStatus)
	}

}

func (pool *Pool) playGame() {
	player1Hand := pool.board[1].hand
	player2Hand := pool.board[2].hand

	winnerId := PlayGame(player1Hand, player2Hand)

	if winnerId == 0 {
		pool.notifyAllPlayers("It's a Tie !!")
	} else {
		pool.notifyWinnerAndLosers(winnerId)
	}

	pool.resetHands()
}

func (pool *Pool) notifyAllPlayers(message string) {
	for c := range pool.Clients {
		c.Conn.WriteJSON(Message{
			MessageType: GM,
			Message:     message,
		})
	}
}

func (pool *Pool) notifyWinnerAndLosers(winnerId int) {
	for c := range pool.Clients {
		if winnerId == c.ID {
			c.Conn.WriteJSON(Message{
				MessageType: GM,
				Message:     "You Win !",
			})
		} else {
			c.Conn.WriteJSON(Message{
				MessageType: GM,
				Message:     "You Lose !",
			})
		}
	}
}

func (pool *Pool) resetHands() {
	pool.board[1].hand = "X"
	pool.board[2].hand = "X"
}

func (pool *Pool) notifyWaitingPlayers(gameStatus int) {
	for client := range pool.Clients {
		if client.ID == gameStatus {
			client.Conn.WriteJSON(Message{
				MessageType: GM,
				Message:     "Waiting for other player",
			})
		} else {
			client.Conn.WriteJSON(Message{
				MessageType: GM,
				Message:     fmt.Sprintf("Waiting for player %d ", client.ID),
			})
		}
	}
}
