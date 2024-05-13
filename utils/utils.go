package utils

import "RockPaperScissor/types"

func SendMessage(messageType int, message string, score int) types.Message {
	return types.Message{
		MessageType: messageType,
		Message:     message,
		Score:       score,
	}
}
