package models

type Game struct {
	Id        string
	PlayerOne *bool
	PlayerTwo *bool
	Status    *bool
	Score     [6]int
	Board     [2]int
	Hands     [2]string
}

type Dict map[int]interface {
}
