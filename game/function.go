package function

import (
	"RockPaper/models"
	"RockPaper/utils"
)

func AssignPlayers(gameToken string, Games map[string]models.Game) bool {
	val, ok := Games[gameToken]

	if !ok {
		Games[gameToken] = models.Game{
			Id:        gameToken,
			PlayerOne: utils.BoolPointer(true),
		}
		return false
	} else {
		val.PlayerTwo = utils.BoolPointer(true)
		val.Status = utils.BoolPointer(true)
		Games[gameToken] = val
		return true
	}
}

// R:1 , P:2 , S:3
// this will check who win then assign in the g.Score with 1 for p1 and 2 for p2
// if both players assigned there hands we return 1 and 0 if they have not yet
func RuleScore(game models.Game) int {
	var ScoreBoard int = get_sum(game.Board)
	g := game
	if ScoreBoard != 0 {
		if ScoreBoard == 5 {
			p1 := g.Hands[1]
			if p1 == "S" {
				g.Score[0] = 1
				return 1
			}
			g.Score[0] = 2
			return 1
		}
	}
	return 0
}

func get_sum(scores [2]int) int {
	var sum int
	for i := range scores {
		sum += scores[i]
	}
	return sum
}
