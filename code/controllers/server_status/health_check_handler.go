package server_status

import (
    "fmt"
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
    _, _ = fmt.Fprintln(httpResponseWriter, "ok")
}

