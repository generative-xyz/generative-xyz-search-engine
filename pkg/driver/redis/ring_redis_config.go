package redis

import (
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/utils"
	"hash/crc32"
	"time"

	"github.com/golang/groupcache/consistenthash"

	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// RingConnection -- redis connection
type RingConnection struct {
	// network    string
	addrs      map[string]string
	password   string
	db         int
	maxRetries int
	poolSize   int
}

// BuildClient -- build single redis client
func (conn *RingConnection) BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error) {
	if len(conn.addrs) == 0 {
		return nil, ErrorMissingRedisAddress
	}

	logger.AtLog.Infof("[redis] shards - addrs: %v, pass: %v, db: %v, pollSize: %v",
		conn.addrs, utils.CensorString(conn.password), conn.db, conn.poolSize)

	if conn.maxRetries <= 0 {
		conn.maxRetries = DefaultMaxRetries
	}
	redisdb := newRing(
		&redis.RingOptions{
			Addrs:       conn.addrs,
			Password:    conn.password, // no password set
			DB:          conn.db,       // use default DB
			PoolSize:    conn.poolSize,
			MaxRetries:  conn.maxRetries,
			PoolTimeout: time.Second * 4,
			NewConsistentHash: func(shards []string) redis.ConsistentHash {
				ch := consistenthash.New(128, crc32.ChecksumIEEE)
				ch.Add(shards...)
				return ch
			},
		},
		opts...,
	)
	return redisdb, nil
}

func newRing(opt *redis.RingOptions, opts ...redistrace.ClientOption) redis.UniversalClient {
	client := redis.NewRing(opt)
	redistrace.WrapClient(client, opts...)
	return client
}
