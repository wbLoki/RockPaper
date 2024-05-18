package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientUrl     string
	RedisAddr     string
	RedisPassword string
	RedisDb       int
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	redisDb, _ := strconv.Atoi(GetEnv("CLIENT_URL", "0"))
	return Config{
		ClientUrl:     GetEnv("CLIENT_URL", "http://localhost:3000/"),
		RedisAddr:     GetEnv("RedisAddr", "localhost:6379"),
		RedisPassword: GetEnv("RedisPassword", "password"),
		RedisDb:       redisDb,
	}
}

func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv("key"); ok {
		return value
	}
	return defaultValue
}
