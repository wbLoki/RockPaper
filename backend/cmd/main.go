package main

import (
	"RockPaperScissor/cmd/api"
	"RockPaperScissor/config"
	"RockPaperScissor/db"
	"RockPaperScissor/pkg"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := db.NewRedisStorage(redis.Options{
		Addr:     config.Envs.RedisAddr,
		Password: config.Envs.RedisPassword,
		DB:       config.Envs.RedisDb,
	})
	pool := pkg.NewPool(rdb)

	go pool.Start()

	server := api.NewApiServer(":8080", pool, rdb)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
