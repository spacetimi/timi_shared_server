package redis_adaptor

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spacetimi/timi_shared_server/code/config"
)

var Client *redis.Client

func Initialize() {
	fmt.Println("Trying to connect to redis at: " + config.GetEnvironmentConfiguration().SharedRedisURL)

	Client = redis.NewClient(&redis.Options {
		Addr:     config.GetEnvironmentConfiguration().SharedRedisURL,
		Password: config.GetEnvironmentConfiguration().SharedRedisPasswd,
		DB:       0,  // use default DB
	})
}

func Ping() bool {
	pong, err := Client.Ping().Result()
	if err != nil {
		panic("Redis ping failed: " + err.Error())
		return false
	}

	fmt.Println(pong)
	return true
}

func Read(key string) string {
	val, err := Client.Get(key).Result()
	if err != nil {
		panic("Failed to find value for key: " + err.Error())
	}
	fmt.Println(key + ":" + val)

	return val
}

func Write(key string, value string) {
	err := Client.Set(key, value, 0).Err()
	if err != nil {
		panic("Failed to set value for key: " + key + ". " + err.Error())
	}
}
