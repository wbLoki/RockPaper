package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientUrl string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		ClientUrl: GetEnv("CLIENT_URL", "http://localhost:3000/"),
	}
}

func GetEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv("key"); ok {
		return value
	}
	return defaultValue
}
