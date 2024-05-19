package types

// Message Formate {"type":2, "message": "paper", "score": 0}
type Message struct {
	ClientId    string `json:"clientId"`
	GameId      string `json:"gameId"`
	MessageType int    `json:"type"`
	Message     string `json:"message"`
	Score       int    `json:"score"`
	Player      Player `json:"player"`
}

type Player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type RedisGame struct {
	ID    string            `json:"id"`
	Lobby []string          `json:"lobby"`
	Hands map[string]string `json:"hands"`
}

type PlayerRedis struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	GameId string `json:"gameId"`
	Score  int    `json:"score"`
	Hand   string `json:"hand"`
}
type UpdatePlayerPayload struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
