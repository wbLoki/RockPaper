package models

type Game struct {
	Id        string
	PlayerOne *bool
	PlayerTwo *bool
	Status    *bool
	Board     [6]int
}
