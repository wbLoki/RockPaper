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
