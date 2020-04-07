package server_status

import (
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/controllers/shared_routes"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
    "strconv"
)

type ConfigHandler struct {     // Implements IRouteHandler
}

func (ch *ConfigHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute(shared_routes.Config, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ch *ConfigHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    _, _ = fmt.Fprintln(httpResponseWriter, "AppName: " + config.GetAppName())
    _, _ = fmt.Fprintln(httpResponseWriter, "Port: " + strconv.Itoa(config.GetEnvironmentConfiguration().Port))

    envString := "unknown"
    switch config.GetEnvironmentConfiguration().AppEnvironment {
    case config.LOCAL: envString = "Local"
    case config.TEST: envString = "Test"
    case config.STAGING: envString = "Staging"
    case config.PRODUCTION: envString = "Production"
    }
    _, _ = fmt.Fprintln(httpResponseWriter, "Environment: " + envString)
    _, _ = fmt.Fprintln(httpResponseWriter, "")

    _, _ = fmt.Fprintln(httpResponseWriter, "Shared MongoDB URL: " + config.GetEnvironmentConfiguration().SharedMongoURL)
    _, _ = fmt.Fprintln(httpResponseWriter, "App MongoDB URL: " + config.GetEnvironmentConfiguration().AppMongoURL)
    _, _ = fmt.Fprintln(httpResponseWriter, "")

    _, _ = fmt.Fprintln(httpResponseWriter, "Shared DB Name: " + config.GetEnvironmentConfiguration().SharedDatabaseName)
    _, _ = fmt.Fprintln(httpResponseWriter, "App DB Name: " + config.GetEnvironmentConfiguration().AppDatabaseName)
    _, _ = fmt.Fprintln(httpResponseWriter, "")

    _, _ = fmt.Fprintln(httpResponseWriter, "Shared Redis URL: " + config.GetEnvironmentConfiguration().SharedRedisURL)
    _, _ = fmt.Fprintln(httpResponseWriter, "App Redis URL: " + config.GetEnvironmentConfiguration().AppRedisURL)
}

