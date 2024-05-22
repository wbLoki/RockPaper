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
	redisDb, _ := strconv.Atoi(GetEnv("REDIS_DB", "0"))
	return Config{
		ClientUrl:     GetEnv("CLIENT_URL", "http://localhost:3000/"),
		RedisAddr:     GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", "password"),
		RedisDb:       redisDb,
	}
}

func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
