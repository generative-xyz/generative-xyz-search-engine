package redis

import (
	"time"

	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// SentinelConnection -- redis connection
type SentinelConnection struct {
	masterGroup       string
	sentinelAddresses []string
	password          string
	db                int
	poolSize          int
}

// BuildClient --
func (conn *SentinelConnection) BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error) {
	if len(conn.sentinelAddresses) == 0 {
		return nil, ErrorMissingRedisAddress
	}

	masterGroup := conn.masterGroup
	if masterGroup == "" {
		masterGroup = "master"
	}

	redisdb := newFailoverClient(
		&redis.FailoverOptions{
			MasterName:    masterGroup,
			SentinelAddrs: conn.sentinelAddresses,
			Password:      conn.password,
			DB:            conn.db,
			PoolSize:      conn.poolSize,
			PoolTimeout:   time.Second * 4,
		},
		opts...)

	return redisdb, nil
}

func newFailoverClient(opt *redis.FailoverOptions, opts ...redistrace.ClientOption) redis.UniversalClient {
	client := redis.NewFailoverClient(opt)
	redistrace.WrapClient(client, opts...)
	return client
}
