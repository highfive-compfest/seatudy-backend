package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     Env.RedisAddress,
		Password: Env.RedisPassword,
		DB:       Env.RedisDatabase,
	})

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("fail to connect redis", err)
	}
	log.Println("connected to redis", ping)

	return client
}