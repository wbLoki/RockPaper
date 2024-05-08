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

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.register:
			p.handleClientRegistration(client)
		case message := <-p.broadcast:
			p.broadcastMessage(message)
		case client := <-p.unregister:
			p.unregisterClient(client)
		case gameStatus := <-p.gameStatus:
			p.handleGameStatus(gameStatus)
		}
	}
}

func (p *Pool) handleClientRegistration(client *Client) {
	client.Conn.WriteJSON(Message{
		MessageType: GM,
		Message:     "Welcome Player " + strconv.Itoa(client.ID),
	})

	for _client := range p.Clients {
		_client.Conn.WriteJSON(Message{
			MessageType: GM,
			Message:     "Player " + strconv.Itoa(client.ID) + " Joined",
		})
	}
	p.Clients[client] = true
	p.board[client.ID] = &Hand{
		client: client,
		hand:   "X",
	}
}

func (p *Pool) broadcastMessage(message Message) {
	for client := range p.Clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (p *Pool) unregisterClient(client *Client) {
	delete(p.Clients, client)
	delete(p.board, client.ID)
}

func (p *Pool) handleGameStatus(gameStatus int) {
	isReady := IsPlayersReady(p)
	if isReady {
		p.playGame()
	} else {
		p.notifyWaitingPlayers(gameStatus)
	}

}

func (p *Pool) playGame() {
	player1Hand := p.board[1].hand
	player2Hand := p.board[2].hand

	winnerId := PlayGame(player1Hand, player2Hand)

	if winnerId == 0 {
		p.notifyAllPlayers("It's a Tie !!")
	} else {
		p.notifyWinnerAndLosers(winnerId)
	}

	p.resetHands()
}

func (p *Pool) notifyAllPlayers(message string) {
	for c := range p.Clients {
		c.Conn.WriteJSON(Message{
			MessageType: GM,
			Message:     message,
		})
	}
}

func (p *Pool) notifyWinnerAndLosers(winnerId int) {
	for c := range p.Clients {
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

func (p *Pool) resetHands() {
	p.board[1].hand = "X"
	p.board[2].hand = "X"
}

func (p *Pool) notifyWaitingPlayers(gameStatus int) {
	for client := range p.Clients {
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
