package pkg

import "RockPaperScissor/types"

type Hand struct {
	client *Client
	hand   string
}

func PlayGame(hand1, hand2 string) int {
	if hand1 == hand2 {
		return 0
	} else if (hand1 == "rock" && hand2 == "scissors") || (hand1 == "paper" && hand2 == "rock") || (hand1 == "scissors" && hand2 == "paper") {
		return 1
	} else {
		return 2
	}
}

func IsPlayersReady(RedisGame types.RedisGame) bool {

	if len(RedisGame.Lobby) == 2 {
		for _, hand := range RedisGame.Hands {
			if hand == "X" {
				return false
			}
		}
	}
	return len(RedisGame.Lobby) == 2
}
