package redis

import (
	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// UnknownConnection --
type UnknownConnection struct{}

// BuildClient --
func (conn *UnknownConnection) BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error) {
	return nil, ErrorRedisClientNotSupported
}
