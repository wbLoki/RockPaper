package pkg

func PlayGame(hand1, hand2 string) int {
	if hand1 == hand2 {
		return 0
	} else if (hand1 == "rock" && hand2 == "scissors") || (hand1 == "paper" && hand2 == "rock") || (hand1 == "scissors" && hand2 == "paper") {
		return 1
	} else {
		return 2
	}
}

func IsPlayersReady(pool *Pool) bool {
	if len(pool.board) == 2 {
		for _, hand := range pool.board {
			if hand.hand == "X" {
				return false
			}
		}
	}
	return len(pool.board) == 2
}
