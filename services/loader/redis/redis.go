package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/vladimir3322/stonent_go/services/loader/config"
)

var ctx = context.Background()

func PushEvent(data []byte) {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisUrl,
	})

	client.RPush(ctx, config.RedisJobQueue, data)
}
