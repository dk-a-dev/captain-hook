package main

import (
	"context"
	"log"
	"os"

	"dkadev.xyz/captainhook/queue"
	redisClient "dkadev.xyz/captainhook/redis"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})

	webhookQueue := make(chan redisClient.WebHookPayload, 100)

	go queue.ProcessHooks(ctx, webhookQueue)
	
	err := redisClient.Subscribe(ctx, client, webhookQueue)
	if err != nil {
		log.Println("Error:", err)
	}

	select {}

}
