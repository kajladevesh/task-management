package main

import (
	"context"
	"notification_service/src/internal/config"
	"notification_service/src/internal/interfaces/subscriber"
	"notification_service/src/internal/usecase"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	ctx := context.Background()
	uc := usecase.NewNotificationUseCase()
	sub := subscriber.NewSubscriber(rdb, uc)

	go sub.Subscribe(ctx, "task_notifications")

	select {} // block forever
}
