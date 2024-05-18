package types

// Message Formate {"type":2, "message": "paper", "score": 0}
type Message struct {
	MessageType int    `json:"type"`
	Message     string `json:"message"`
	Score       int    `json:"score"`
}

type RedisGame struct {
	ID    string            `json:"id"`
	Lobby []string          `json:"lobby"`
	Hands map[string]string `json:"hands"`
}

type PlayerRedis struct {
	Name   string `json:"name"`
	GameId string `json:"gameId"`
	Score  int    `json:"score"`
	Hand   string `json:"hand"`
}
