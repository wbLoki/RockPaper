package utils

import (
	"RockPaperScissor/types"
	"context"
	"encoding/json"
	"math/rand"

	"github.com/redis/go-redis/v9"
)

func SendMessage(messageType int, message string, score int) types.Message {
	return types.Message{
		MessageType: messageType,
		Message:     message,
		Score:       score,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetFromRedis(rdb *redis.Client, key string, v interface{}) error {
	var ctx = context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), v); err != nil {
		return err
	}

	return nil
}

func SetRedis(rdb *redis.Client, key string, v any) error {

	val, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var ctx = context.Background()
	if err := rdb.Set(ctx, key, string(val), 0).Err(); err != nil {
		return err
	}

	return nil
}

func removeItem(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}

func RemoveItemByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return removeItem(slice, i)
		}
	}
	return slice
}
