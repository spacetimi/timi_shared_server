package server_status

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/code/controllers/shared_routes"
	"github.com/spacetimi/timi_shared_server/v2/code/core/controller"
)

type ConfigHandler struct { // Implements IRouteHandler
}

func (ch *ConfigHandler) Routes() []controller.Route {
	return []controller.Route{
		controller.NewRoute(shared_routes.Config, []controller.RequestMethodType{controller.GET, controller.POST}),
	}
}

func (ch *ConfigHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	if config.GetEnvironmentConfiguration().AppEnvironment == config.PRODUCTION {
		// TODO: Figure out how to enable this on production and just rpevent access to this endpoint if its out of vpn
		_, _ = fmt.Fprintln(httpResponseWriter, "not on production")
		return
	}

	_, _ = fmt.Fprintln(httpResponseWriter, "AppName: "+config.GetAppName())
	_, _ = fmt.Fprintln(httpResponseWriter, "Port: "+strconv.Itoa(config.GetEnvironmentConfiguration().Port))

	envString := "unknown"
	switch config.GetEnvironmentConfiguration().AppEnvironment {
	case config.LOCAL:
		envString = "Local"
	case config.TEST:
		envString = "Test"
	case config.STAGING:
		envString = "Staging"
	case config.PRODUCTION:
		envString = "Production"
	}
	_, _ = fmt.Fprintln(httpResponseWriter, "Environment: "+envString)
	_, _ = fmt.Fprintln(httpResponseWriter, "")

	// TODO: Also print the hash of the latest commit in app and shared?

	// TODO: Figure out how to print these configs (probably by hanging on to the Config object in each service and returning that)
	//_, _ = fmt.Fprintln(httpResponseWriter, "Shared MongoDB URL: " + config.GetEnvironmentConfiguration().SharedMongoURL)
	//_, _ = fmt.Fprintln(httpResponseWriter, "App MongoDB URL: " + config.GetEnvironmentConfiguration().AppMongoURL)
	//_, _ = fmt.Fprintln(httpResponseWriter, "")
	//
	//_, _ = fmt.Fprintln(httpResponseWriter, "Shared DB Name: " + config.GetEnvironmentConfiguration().SharedDatabaseName)
	//_, _ = fmt.Fprintln(httpResponseWriter, "App DB Name: " + config.GetEnvironmentConfiguration().AppDatabaseName)
	//_, _ = fmt.Fprintln(httpResponseWriter, "")
	//
	//_, _ = fmt.Fprintln(httpResponseWriter, "Shared Redis URL: " + config.GetEnvironmentConfiguration().SharedRedisURL)
	//_, _ = fmt.Fprintln(httpResponseWriter, "App Redis URL: " + config.GetEnvironmentConfiguration().AppRedisURL)
}
