package server_status

import (
    "errors"
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/core/adaptors/redis_adaptor"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
)

type HealthCheckHandler struct {     // Implements IRouteHandler
}

func (hch *HealthCheckHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute("/healthCheck", []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (hch *HealthCheckHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    ok, err := performHealthChecks()

    if ok {
        _, _ = fmt.Fprintln(httpResponseWriter, "ok")
    } else {
        _, _ = fmt.Fprintln(httpResponseWriter, err.Error())
    }
}

func performHealthChecks() (bool, error) {
    _, err := redis_adaptor.Ping()
    if err != nil {
        return false, errors.New("check redis failed with error: " + err.Error())
    }

    // TODO: Add health check for mongo_adaptor
    // TODO: Add health check for metadata_service
    // TODO: Add health check for storage_service
    // TODO: Add health check for identity_service

    return true, nil
}

