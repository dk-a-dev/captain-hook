package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type WebHookPayload struct {
	Url       string `json:"url"`
	WebhookId string `json:"webhookId"`
	Data      struct {
		Id      string `json:"id"`
		Payment string `json:"payment"`
		Event   string `json:"event"`
		Date    string `json:"created"`
	} `json:"data"`
}

func Subscribe(ctx context.Context, client *redis.Client, webhookQueue chan WebHookPayload) error {
	pubSub := client.Subscribe(ctx, "payment")

	defer func(pubSub *redis.PubSub) {
		if err := pubSub.Close(); err != nil {
			log.Println("Error closing PubSub:", err)
		}
	}(pubSub)

	var payload WebHookPayload
	for {
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Println("Error un-marshalling payload:", err)
			continue
		}

		webhookQueue <- payload
	}
}
