package subscriber

import (
	"context"
	"log"
	"notification_service/src/internal/usecase"

	"github.com/go-redis/redis/v8"
)

type Subscriber struct {
	rdb     *redis.Client
	usecase *usecase.NotificationUseCase
}

func NewSubscriber(rdb *redis.Client, uc *usecase.NotificationUseCase) *Subscriber {
	return &Subscriber{rdb: rdb, usecase: uc}
}

func (s *Subscriber) Subscribe(ctx context.Context, channel string) {
	pubsub := s.rdb.Subscribe(ctx, channel)

	log.Printf("Subscribed to Redis channel: %s", channel)

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Message received on channel %s: %s", channel, msg.Payload)
		err := s.usecase.HandleNotification(ctx, msg.Payload)
		if err != nil {
			log.Printf("Error handling notification: %v", err)
		}
	}
}
