package redis

import (
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/utils"

	"time"

	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

const (
	// DefaultMaxRetries --
	DefaultMaxRetries int = 3
)

// SingleConnection -- redis connection
type SingleConnection struct {
	// network    string
	address    string
	password   string
	db         int
	maxRetries int
	poolSize   int
}

// BuildClient -- build single redis client
func (conn *SingleConnection) BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error) {
	if conn.address == "" {
		return nil, ErrorMissingRedisAddress
	}

	logger.AtLog.Infof("[redis] single - address: %v, pass: %v, db: %v, pollSize: %v",
		conn.address, utils.CensorString(conn.password), conn.db, conn.poolSize)

	if conn.maxRetries <= 0 {
		conn.maxRetries = DefaultMaxRetries
	}
	redisdb := redistrace.NewClient(
		&redis.Options{
			Addr:        conn.address,
			Password:    conn.password, // no password set
			DB:          conn.db,       // use default DB
			PoolSize:    conn.poolSize,
			MaxRetries:  conn.maxRetries,
			PoolTimeout: time.Second * 4,
		},
		opts...,
	)
	return redisdb, nil
}
