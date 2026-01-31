package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisPublisher struct {
	Client *redis.Client
}

func NewRedisPublisher(addr, password string) *RedisPublisher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisPublisher{Client: rdb}
}

func (p *RedisPublisher) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.Client.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	//log.Printf("Published to channel %s: %s", channel, data)
	return nil
}
