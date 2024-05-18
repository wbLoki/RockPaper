package pkg

import (
	"RockPaperScissor/types"
	"RockPaperScissor/utils"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Pool struct {
	Clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan types.Message
	board      map[int]*Hand
	gameId     chan string
	rdb        *redis.Client
}

func NewPool(rdb *redis.Client) *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan types.Message),
		board:      make(map[int]*Hand),
		gameId:     make(chan string),
		rdb:        rdb,
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
		case gameId := <-p.gameId:
			p.handleGameStatus(gameId)
		}
	}
}

func (p *Pool) handleClientRegistration(client *Client) {
	client.Conn.WriteJSON(utils.SendMessage(Chat, "Welcome Player "+client.ID, client.gameBoard.score))
	for _client := range p.Clients {
		_client.Conn.WriteJSON(
			utils.SendMessage(GM, "Player "+client.ID+" Joined", _client.gameBoard.score))
	}
	p.Clients[client] = true

}

func (p *Pool) broadcastMessage(message types.Message) {
	for client := range p.Clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (p *Pool) unregisterClient(client *Client) {
	delete(p.Clients, client)
}

func (p *Pool) handleGameStatus(gameId string) {
	var RedisGame types.RedisGame
	utils.GetFromRedis(p.rdb, gameId, &RedisGame)
	isReady := IsPlayersReady(RedisGame)
	if isReady {
		p.playGame(RedisGame)
	} else {
		p.notifyWaitingPlayers()
	}

}

func (p *Pool) playGame(RedisGame types.RedisGame) {
	player1Hand := RedisGame.Hands[RedisGame.Lobby[0]]
	player2Hand := RedisGame.Hands[RedisGame.Lobby[1]]

	winnerId := PlayGame(player1Hand, player2Hand)

	if winnerId == 0 {
		p.notifyAllPlayers("It's a Tie !!")
	} else {
		p.notifyWinnerAndLosers(winnerId-1, RedisGame)
	}

	p.resetHands(RedisGame)
}

func (p *Pool) notifyAllPlayers(message string) {
	for c := range p.Clients {
		c.Conn.WriteJSON(utils.SendMessage(GM, message, c.gameBoard.score))
	}
}

func (p *Pool) notifyWinnerAndLosers(winnerId int, RedisGame types.RedisGame) {
	wPlayer := RedisGame.Lobby[winnerId]
	for c := range p.Clients {
		if wPlayer == c.ID {
			c.gameBoard.score++
			c.Conn.WriteJSON(utils.SendMessage(GM, "You Win !", c.gameBoard.score))
		} else {
			c.Conn.WriteJSON(utils.SendMessage(GM, "You Lose !", c.gameBoard.score))
		}
	}
}

func (p *Pool) resetHands(RedisGame types.RedisGame) {
	RedisGame.Hands[RedisGame.Lobby[0]] = "X"
	RedisGame.Hands[RedisGame.Lobby[1]] = "X"
	if err := utils.SetRedis(p.rdb, RedisGame.ID, RedisGame); err != nil {
		log.Fatal(err)
	}
}

func (p *Pool) notifyWaitingPlayers() {
	for client := range p.Clients {
		var message string
		message = "Waiting for player"
		client.Conn.WriteJSON(utils.SendMessage(GM, message, client.gameBoard.score))
	}
}
