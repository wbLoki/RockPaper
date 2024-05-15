package types

// Message Formate {"type":2, "message": "paper", "score": 0}
type Message struct {
	MessageType int    `json:"type"`
	Message     string `json:"message"`
	Score       int    `json:"score"`
}
