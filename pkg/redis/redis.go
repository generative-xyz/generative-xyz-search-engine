package redis

import (
	"context"
	"generative-xyz-search-engine/pkg/driver/redis"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Client interface {
	Client() *redis.AtRedis
	Get(ctx context.Context, key string, result interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Dels(ctx context.Context, keys ...string) error
	DelPrefix(ctx context.Context, prefix string) error
	HGet(ctx context.Context, key, field string, result interface{}) error
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key string) error
}

type clientImpl struct {
	client *redis.AtRedis
}

func NewClient() Client {
	return &clientImpl{redis.CreateRedisConnection(nil)}
}

func (s *clientImpl) Client() *redis.AtRedis {
	return s.client
}

func (s *clientImpl) Get(ctx context.Context, key string, result interface{}) error {
	return s.client.Get(ctx, key).Scan(result)
}

func (s *clientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *clientImpl) Del(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *clientImpl) Dels(ctx context.Context, keys ...string) error {
	if viper.GetString("REDIS_CLIENT_TYPE") == redis.Ring {
		val, ok := s.client.UniversalClient.(*redisv8.Ring)
		if ok {
			if err := val.ForEachShard(ctx, func(ctx context.Context, client *redisv8.Client) error {
				return client.Del(ctx, keys...).Err()
			}); err != nil {
				return err
			}
			return nil
		}
	}
	return s.client.Del(ctx, keys...).Err()
}

func (s *clientImpl) DelPrefix(ctx context.Context, prefix string) error {
	if viper.GetString("REDIS_CLIENT_TYPE") == redis.Ring {
		val, ok := s.client.UniversalClient.(*redisv8.Ring)
		if ok {
			if err := val.ForEachShard(ctx, func(ctx context.Context, client *redisv8.Client) error {
				return s.delKeys(ctx, prefix, client)
			}); err != nil {
				return err
			}
			return nil
		}
	}
	return s.delKeys(ctx, prefix, s.client.UniversalClient)
}

func (s *clientImpl) delKeys(ctx context.Context, prefix string, client redisv8.UniversalClient) error {
	keys, err := client.Keys(ctx, prefix+"*").Result()
	if err == nil && len(keys) > 0 {
		if err := client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *clientImpl) HGet(ctx context.Context, key, field string, result interface{}) error {
	return s.client.HGet(ctx, key, field).Scan(result)
}

func (s *clientImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return s.client.HSet(ctx, key, values).Err()
}

func (s *clientImpl) HDel(ctx context.Context, key string) error {
	return s.client.HDel(ctx, key).Err()
}
