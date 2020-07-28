package server_status

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spacetimi/timi_shared_server/v2/code/controllers/shared_routes"
	"github.com/spacetimi/timi_shared_server/v2/code/core/adaptors/redis_adaptor"
	"github.com/spacetimi/timi_shared_server/v2/code/core/controller"
)

type HealthCheckHandler struct { // Implements IRouteHandler
}

func (hch *HealthCheckHandler) Routes() []controller.Route {
	return []controller.Route{
		controller.NewRoute(shared_routes.HealthCheck, []controller.RequestMethodType{controller.GET, controller.POST}),
	}
}

func (hch *HealthCheckHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	ok, err := performHealthChecks(request.Context())

	if ok {
		_, _ = fmt.Fprintln(httpResponseWriter, "ok")
	} else {
		_, _ = fmt.Fprintln(httpResponseWriter, err.Error())
	}
}

func performHealthChecks(ctx context.Context) (bool, error) {
	// TODO: Check if each service is being used before checking their health

	_, err := redis_adaptor.Ping(ctx)
	if err != nil {
		return false, errors.New("check redis failed with error: " + err.Error())
	}

	// TODO: Add health check for mongo_adaptor
	// TODO: Add health check for metadata_service
	// TODO: Add health check for storage_service
	// TODO: Add health check for identity_service

	return true, nil
}
