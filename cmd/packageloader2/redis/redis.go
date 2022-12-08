package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Add(key string, value interface{}) error {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (interface{}, error) {
	ret, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ret, err
	}
	return ret, nil
}
