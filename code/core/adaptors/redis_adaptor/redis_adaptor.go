package redis_adaptor

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/spacetimi/timi_shared_server/code/config"
	"time"
)

var _client *redis.Client

var EXPIRATION_DEFAULT time.Duration = 48 * time.Hour

func Initialize() {
	_client = redis.NewClient(&redis.Options {
		Addr:     config.GetEnvironmentConfiguration().SharedRedisURL,
		Password: config.GetEnvironmentConfiguration().SharedRedisPasswd,
		DB:       0,  // use default DB
	})
}

func Ping(ctx context.Context) (bool, error) {
	_, err := _client.Ping(ctx).Result()
	if err != nil {
		return false, errors.New("error pinging redis: " + err.Error())
	}

	return true, nil
}

func Read(key string, ctx context.Context) (string, bool) {
	val, err := _client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}

func Write(key string, value string, expiration time.Duration, ctx context.Context) error {
	err := _client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.New("error writing value for key: " + err.Error())
	}

	return nil
}

func Delete(key string, ctx context.Context) error {
	err := _client.Del(ctx, key).Err()
	if err != nil {
		return errors.New("error deleting key: " + err.Error())
	}

	return nil
}
