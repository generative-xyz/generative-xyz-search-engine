package redis

import (
	"errors"
	"generative-xyz-search-engine/utils"
	"strings"

	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Connection --
type Connection interface {
	BuildClient(opts ...redistrace.ClientOption) (redis.UniversalClient, error)
}

var (
	// ErrorMissingRedisAddress --
	ErrorMissingRedisAddress = errors.New("missing redis address")

	// ErrorRedisClientNotSupported --
	ErrorRedisClientNotSupported = errors.New("redis client not supported")
)

const (
	// Sentinel type
	Sentinel = "sentinel"
	// Cluster type
	Cluster = "cluster"
	// Single type
	Single = "single"
	// Ring type
	Ring = "ring"

	// DefaultPoolSize --
	DefaultPoolSize = 100
)

// DefaultRedisConnectionFromConfig -- load connection settings in config with default key
func DefaultRedisConnectionFromConfig() Connection {
	redisClientType := viper.GetString("REDIS_CLIENT_TYPE")
	if utils.IsStringEmpty(redisClientType) {
		redisClientType = Single
	}
	poolSize := viper.GetInt("REDIS_POOL_SIZE")
	if poolSize <= 0 {
		poolSize = DefaultPoolSize
	}

	switch redisClientType {
	case Sentinel:
		return &SentinelConnection{
			masterGroup:       viper.GetString("REDIS_SENTINEL_MASTER"),
			sentinelAddresses: viper.GetStringSlice("REDIS_SENTINEL_ADDRESS"),
			password:          viper.GetString("REDIS_PASSWORD"),
			db:                viper.GetInt("REDIS_DATABASE"),
			poolSize:          poolSize,
		}
	case Cluster:
		return &ClusterConnection{
			clusterAddresses: viper.GetStringSlice("REDIS_CLUSTER_ADDRESS"),
			password:         viper.GetString("REDIS_PASSWORD"),
			poolSize:         poolSize,
		}
	case Single:
		return &SingleConnection{
			address:    viper.GetString("REDIS_ADDRESS"),
			password:   viper.GetString("REDIS_PASSWORD"),
			db:         viper.GetInt("REDIS_DATABASE"),
			maxRetries: viper.GetInt("REDIS_MAX_RETRIES"),
			poolSize:   poolSize,
		}
	case Ring:
		addrs := map[string]string{}
		shards := make([]string, 0)
		shardsStr := viper.GetString("REDIS_SHARDS")
		if strings.Contains(shardsStr, ",") {
			shards = strings.Split(shardsStr, ",")
		} else {
			shards = viper.GetStringSlice("REDIS_SHARDS")
		}
		for _, host := range shards {
			host = strings.TrimSpace(host)
			addrs[host] = host
		}
		return &RingConnection{
			addrs:      addrs,
			password:   viper.GetString("REDIS_PASSWORD"),
			db:         viper.GetInt("REDIS_DATABASE"),
			maxRetries: viper.GetInt("REDIS_MAX_RETRIES"),
			poolSize:   poolSize,
		}

	default:
		return &UnknownConnection{}
	}
}

// NewRedisConfig --
func NewRedisConfig(add string, db int) Connection {
	return &SingleConnection{
		address:  add,
		db:       db,
		poolSize: DefaultPoolSize,
	}
}

// NewRedisConfigWithPool --
func NewRedisConfigWithPool(add string, db, poolSize int) Connection {
	return &SingleConnection{
		address:  add,
		db:       db,
		poolSize: poolSize,
	}
}
