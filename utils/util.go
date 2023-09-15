package utils

import "github.com/gorilla/websocket"

func BoolPointer(s bool) *bool {
	return &s
}

func CastMessage(message []byte, mType int, connections []*websocket.Conn) {
	for _, connection := range connections {
		connection.WriteMessage(mType, message)
	}
}
