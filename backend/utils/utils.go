package utils

import (
	"RockPaperScissor/types"
	"math/rand"
)

func SendMessage(messageType int, message string, score int) types.Message {
	return types.Message{
		MessageType: messageType,
		Message:     message,
		Score:       score,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
