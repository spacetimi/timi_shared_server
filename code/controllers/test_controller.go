package controllers

import (
	"fmt"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func TestController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	// TODO: Print server version, environment, etc
	output := "Environment: "
	switch config.GetEnvironmentConfiguration().AppEnvironment {
	case config.LOCAL: output += "Local"
	case config.TEST: output += "Test"
	case config.STAGING: output += "Staging"
	case config.PRODUCTION: output += "Production"
	}
	fmt.Fprintln(httpResponseWriter, output)

	fmt.Fprintln(httpResponseWriter, "AppName: " + config.GetAppName())
	fmt.Fprintln(httpResponseWriter, "Port: " + strconv.Itoa(config.GetEnvironmentConfiguration().Port))
	fmt.Fprintln(httpResponseWriter, "MongoDB URL: " + config.GetEnvironmentConfiguration().SharedMongoURL)
	fmt.Fprintln(httpResponseWriter, "Redis URL: " + config.GetEnvironmentConfiguration().SharedRedisURL)
	fmt.Fprintln(httpResponseWriter, "Shared DB Name: " + config.GetEnvironmentConfiguration().SharedDatabaseName)
	fmt.Fprintln(httpResponseWriter, "App DB Name: " + config.GetEnvironmentConfiguration().AppDatabaseName)
}

func PingRedisController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	pingOk, err := redis_adaptor.Ping()
	if pingOk {
		fmt.Fprintln(httpResponseWriter, "Redis ping successful")
		redis_adaptor.Write("test_key", strconv.Itoa(rand.Int()))
		value, ok := redis_adaptor.Read("test_key")
		if ok {
			fmt.Fprintln(httpResponseWriter, "Value: " + value)
		} else {
			fmt.Fprintln(httpResponseWriter, "Redis read returned empty")
		}

	} else {
		fmt.Fprintf(httpResponseWriter, "Redis ping failed|error=" + err.Error())
	}
}

func PanicController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	panic("some panic occurred")
}
func FatalController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	log.Fatal("log fatal occurred")
}
