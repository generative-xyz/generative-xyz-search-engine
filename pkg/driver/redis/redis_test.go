package redis

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func countNil(l []interface{}) (count int) {
	for k := range l {
		if l[k] == nil {
			count++
		}
	}
	return
}
func TestRedisGet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	result := client.HMGet(context.Background(), "sakat:38a8f27a0b199cf412a5a3c6613356e5",
		[]string{"shop_id", "user_id", "apps_shops"}...)
	v, e := result.Result()
	if e != nil || len(v) == 0 || countNil(v) == len(v) {
		t.Error("Not found", len(v), e)
		return
	}

	shopID2, ok := v[0].(string)
	if ok {
		fmt.Println("shopID2:", shopID2, reflect.TypeOf(v[0]))
		fmt.Printf("%#v %T", v[0], v[0])
	}

	userID, ok := v[1].(string)
	if ok {
		fmt.Println("userID2:", userID, reflect.TypeOf(v[1]))
	}

	fmt.Println(v...)
	t.Error()

}

func TestRedisCluster(t *testing.T) {
	viper.Set("redis.clientType", "cluster")
	viper.Set("redis.cluster.addresses", []string{"127.0.0.1:30001", "127.0.0.1:30002", "127.0.0.1:30003"})
	conn := DefaultRedisConnectionFromConfig()
	fmt.Printf("%#v %v", conn, conn)
	bkrd, err := NewConnection(conn)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("slot: ", bkrd.Slots)
	t.Error("")
}

func TestRedisRing(t *testing.T) {
	viper.Set("REDIS_CLIENT_TYPE", "ring")
	viper.Set("REDIS_SHARDS", "127.0.0.1:6380, 127.0.0.1:6381")
	conn := DefaultRedisConnectionFromConfig()
	fmt.Printf("%#v %v", conn, conn)
	bkrd, err := NewConnection(conn)
	if err != nil {
		t.Error(err)
		return
	}

	// test operation
	key := "test"
	value := "tada"
	ctx := context.Background()

	bkrd.Set(ctx, key, value, time.Minute*1)
	strResult := bkrd.Get(ctx, key)
	if strResult == nil || strResult.Val() != value {
		t.Errorf("cannot get and set value")
	}

	bkrd.Del(ctx, key)
	strResult = bkrd.Get(ctx, key)
	if strResult == nil || strResult.Val() != "" {
		t.Errorf("can not delete")
	}
}

func TestRedisSAdd(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	key := "cart"
	client.Del(ctx, key)
	client.SAdd(ctx, key, "token A")
	client.SAdd(ctx, key, "token B")
	client.SAdd(ctx, key, "token C")
	client.SAdd(ctx, key, "token B")

	value, err := client.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println("Error: ", err)
		t.Fail()
	}
	fmt.Println("Result: ", value)
}

func TestRedisBLPop(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()
	i := 0
	for {
		i++
		r := client.BLPop(context.Background(), 1*time.Second, "test")
		s, err := r.Result()
		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 03:04:05.99999999"), "-", i, "timed out. ", err)
			continue
		}
		fmt.Println(time.Now().Format("2006-01-02 03:04:05.99999999"), "-", i, ": ", strings.Join(s, ", "))
	}
}
