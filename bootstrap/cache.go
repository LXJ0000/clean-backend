package bootstrap

import (
	"log"
	"time"

	"github.com/LXJ0000/clean-backend/utils/cache"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func NewRedisCache(env *Env) cache.RedisCache {
	cmd := redis.NewClient(&redis.Options{
		Addr: env.RedisAddr,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return cache.NewRedisCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}
