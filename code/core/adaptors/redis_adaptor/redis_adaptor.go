package redis_adaptor

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/spacetimi/timi_shared_server/code/config"
)

var Client *redis.Client

func Initialize() {
	Client = redis.NewClient(&redis.Options {
		Addr:     config.GetEnvironmentConfiguration().SharedRedisURL,
		Password: config.GetEnvironmentConfiguration().SharedRedisPasswd,
		DB:       0,  // use default DB
	})
}

func Ping() (bool, error) {
	_, err := Client.Ping().Result()
	if err != nil {
		return false, errors.New("error pinging redis: " + err.Error())
	}

	return true, nil
}

func Read(key string) (string, bool) {
	val, err := Client.Get(key).Result()
	if err != nil {
		return "", false
	}

	return val, true
}

func Write(key string, value string) error {
	err := Client.Set(key, value, 0).Err()
	if err != nil {
		return errors.New("error writing value for key: " + err.Error())
	}

	return nil
}
