package redis

import (
	"context"
	"fmt"

	"generative-xyz-search-engine/pkg/logger"

	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"

	redis "github.com/go-redis/redis/v8"
)

// CreateRedisConnection --
func CreateRedisConnection(conn Connection, opts ...redistrace.ClientOption) *AtRedis {
	if conn == nil {
		conn = DefaultRedisConnectionFromConfig()
	}
	db, err := NewConnection(conn, opts...)
	if err != nil {
		logger.AtLog.Panic(err)
	}
	return db
}

// AtRedis --
type AtRedis struct {
	redis.UniversalClient
	Slots []redis.ClusterSlot
}

// NewConnection -- open connection to db
func NewConnection(conn Connection, opts ...redistrace.ClientOption) (*AtRedis, error) {
	var err error
	c, err := conn.BuildClient(opts...)
	if err != nil {
		logger.AtLog.Error("[redis] Could not build redis client, details: ", err)
		return nil, err
	}
	pong, err := c.Ping(context.Background()).Result()
	if err != nil {
		logger.AtLog.Error("[redis] Could not ping to redis, details: ", err)
		return nil, err
	}
	logger.AtLog.Info("[redis] Ping to redis: ", pong)
	cs := getClusterInfo(c)
	return &AtRedis{c, cs}, nil
}

func getClusterInfo(c redis.UniversalClient) []redis.ClusterSlot {
	cs := []redis.ClusterSlot{}
	if ci := c.ClusterInfo(context.Background()); ci.Err() == nil {
		csr := c.ClusterSlots(context.Background())
		var err error
		cs, err = csr.Result()
		if err != nil {
			logger.AtLog.Error("[redis] Cannot get cluster slots")
		}
	}
	return cs
}

// NewConnectionFromExistedClient --
func NewConnectionFromExistedClient(c redis.UniversalClient) *AtRedis {
	cs := getClusterInfo(c)
	return &AtRedis{c, cs}
}

// Close -- close connection
func (r *AtRedis) Close() error {
	if r != nil {
		return r.UniversalClient.Close()
	}
	return nil
}

// GetClient --
func (r *AtRedis) GetClient() redis.UniversalClient {
	return r.UniversalClient
}

// GetClusterSlots -
func (r *AtRedis) GetClusterSlots() ([]redis.ClusterSlot, error) {
	res := r.ClusterSlots(context.Background())
	return res.Result()
}

// GetRedisSlot -
func (r *AtRedis) GetRedisSlot(key string) int {
	return Slot(key)
}

// GetRedisSlotID -
func (r *AtRedis) GetRedisSlotID(key string) string {
	return GetSlotID(key, r.Slots)
}

// IsInSlot -
func IsInSlot(key string, slot redis.ClusterSlot) bool {
	s := Slot(key)
	return slot.Start <= s && s <= slot.End
}

// GetSlotID -
func GetSlotID(key string, slots []redis.ClusterSlot) string {
	s := Slot(key)
	for k := range slots {
		slot := slots[k]
		if slot.Start <= s && s <= slot.End {
			return fmt.Sprintf("%v-%v", slot.Start, slot.End)
		}
	}
	return ""
}
