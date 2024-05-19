package pkg

import (
	"RockPaperScissor/types"
	"RockPaperScissor/utils"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Pool struct {
	Clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan types.Message
	gameId     chan string
	rdb        *redis.Client
}

func NewPool(rdb *redis.Client) *Pool {
	return &Pool{
		Clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan types.Message),
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

func (p *Pool) broadcastToLobby(message types.Message) error {
	var RedisGame types.RedisGame
	if err := utils.GetFromRedis(p.rdb, message.GameId, &RedisGame); err != nil {
		fmt.Println(err)
		return err
	}

	for _, clientId := range RedisGame.Lobby {
		p.Clients[clientId].Conn.WriteJSON(message)
	}
	return nil
}

func (p *Pool) handleClientRegistration(client *Client) {

	var RedisGame types.RedisGame
	if err := utils.GetFromRedis(p.rdb, client.GameId, &RedisGame); err != nil {
		log.Fatal(err)
	}

	p.Clients[client.ID] = client

	for _, clientId := range RedisGame.Lobby {
		p.Clients[clientId].Conn.WriteJSON(
			utils.SendMessage(GM, "Player "+client.ID+" Joined", 0))
	}

}

func (p *Pool) broadcastMessage(message types.Message) {
	var RedisGame types.RedisGame
	if err := utils.GetFromRedis(p.rdb, message.GameId, &RedisGame); err != nil {
		fmt.Println(err)
		return
	}

	for _, clientId := range RedisGame.Lobby {
		if _, ok := p.Clients[clientId]; ok {
			p.Clients[clientId].Conn.WriteJSON(message)
		}
	}
}

func (p *Pool) unregisterClient(client *Client) {
	var RedisGame types.RedisGame
	if err := utils.GetFromRedis(p.rdb, client.GameId, &RedisGame); err != nil {
		fmt.Println(err)
	}
	RedisGame.Lobby = utils.RemoveItemByValue(RedisGame.Lobby, client.ID)
	delete(RedisGame.Hands, client.ID)
	delete(p.Clients, client.ID)

	if err := utils.SetRedis(p.rdb, client.GameId, RedisGame); err != nil {
		fmt.Println(err)
	}

}

func (p *Pool) handleGameStatus(gameId string) {
	var RedisGame types.RedisGame
	utils.GetFromRedis(p.rdb, gameId, &RedisGame)
	isReady := IsPlayersReady(RedisGame)
	if isReady {
		p.playGame(RedisGame)
	} else {
		p.broadcastToLobby(types.Message{
			Message:     "Waiting for player",
			MessageType: GM,
			GameId:      gameId})
	}

}

func (p *Pool) playGame(RedisGame types.RedisGame) {
	player1Hand := RedisGame.Hands[RedisGame.Lobby[0]]
	player2Hand := RedisGame.Hands[RedisGame.Lobby[1]]

	winnerId := PlayGame(player1Hand, player2Hand)

	if winnerId == 0 {
		p.broadcastToLobby(types.Message{
			Message:     "It's a Tie !!",
			MessageType: GM,
			GameId:      RedisGame.ID,
		})
	} else {
		p.notifyWinnerAndLosers(winnerId-1, RedisGame)
	}

	p.resetHands(RedisGame)
}

func (p *Pool) notifyWinnerAndLosers(winnerId int, RedisGame types.RedisGame) {
	wPlayer := RedisGame.Lobby[winnerId]
	for _, clientId := range RedisGame.Lobby {
		if wPlayer == clientId {
			var PlayerRedis types.PlayerRedis
			if err := utils.GetFromRedis(p.rdb, clientId, &PlayerRedis); err != nil {
				fmt.Println(err)
			}

			PlayerRedis.Score += 1

			if err := utils.SetRedis(p.rdb, clientId, PlayerRedis); err != nil {
				fmt.Println(err)
			}
			player := types.Player{
				Name: PlayerRedis.Name,
			}
			message := types.Message{
				MessageType: GE,
				ClientId:    clientId,
				GameId:      RedisGame.ID,
				Message:     "You Win !",
				Score:       PlayerRedis.Score,
				Player:      player,
			}
			p.Clients[clientId].Conn.WriteJSON(message)
		} else {
			var PlayerRedis types.PlayerRedis
			if err := utils.GetFromRedis(p.rdb, clientId, &PlayerRedis); err != nil {
				fmt.Println(err)
			}
			player := types.Player{
				Name: PlayerRedis.Name,
			}
			message := types.Message{
				MessageType: GE,
				ClientId:    clientId,
				GameId:      RedisGame.ID,
				Message:     "You Lose !",
				Score:       PlayerRedis.Score,
				Player:      player,
			}
			p.Clients[clientId].Conn.WriteJSON(message)
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
