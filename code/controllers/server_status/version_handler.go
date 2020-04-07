package server_status

import (
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "net/http"
)

type VersionHandler struct {     // Implements IRouteHandler
}

func (vh *VersionHandler) Routes() []controller.Route {
    return []controller.Route {
        controller.NewRoute("/version", []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (vh *VersionHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
    _, _ = fmt.Fprintln(httpResponseWriter, "TODO")
}

