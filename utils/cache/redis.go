package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisCache(cmd redis.Cmdable, expiration time.Duration) RedisCache {
	return RedisCache{cmd: cmd, expiration: expiration}
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.cmd.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.cmd.Get(ctx, key).Result()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisCache) LuaWithreturnInt(ctx context.Context, luaPath string, key []string, args ...interface{}) (int, error) {
	return c.cmd.Eval(ctx, luaPath, key, args).Int()
}

func (c *RedisCache) LuaWithReturnBool(ctx context.Context, luaPath string, key []string, args ...interface{}) (bool, error) {
	return c.cmd.Eval(ctx, luaPath, key, args).Bool()
}

func (c *RedisCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.cmd.HSet(ctx, key, values).Err()
}

func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.cmd.HGetAll(ctx, key).Result()
}
