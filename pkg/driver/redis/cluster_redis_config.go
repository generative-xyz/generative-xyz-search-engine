package redis

import (
	"generative-xyz-search-engine/pkg/logger"
	"time"

	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

// ClusterConnection -- redis connection
type ClusterConnection struct {
	clusterAddresses []string
	password         string
	poolSize         int
}

// BuildClient --
func (conn *ClusterConnection) BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error) {
	if len(conn.clusterAddresses) == 0 {
		return nil, ErrorMissingRedisAddress
	}

	logger.AtLog.Infof("[redis] Create cluster client to %v", conn.clusterAddresses)

	redisdb := newClusterClient(
		&redis.ClusterOptions{
			Addrs:       conn.clusterAddresses,
			Password:    conn.password,
			PoolSize:    conn.poolSize,
			PoolTimeout: time.Second * 4,
		},
		opts...,
	)

	return redisdb, nil
}

func newClusterClient(opt *redis.ClusterOptions, opts ...redistrace.ClientOption) redis.UniversalClient {
	client := redis.NewClusterClient(opt)
	redistrace.WrapClient(client, opts...)
	return client
}
