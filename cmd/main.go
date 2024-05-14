package main

import (
	"RockPaperScissor/cmd/api"
	"RockPaperScissor/pkg"
	"log"
)

func main() {
	hub := pkg.NewHub()

	server := api.NewApiServer(":8080", hub)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
