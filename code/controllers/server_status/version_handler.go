package server_status

import (
	"fmt"
	"net/http"

	"github.com/spacetimi/timi_shared_server/v2/code/controllers/shared_routes"
	"github.com/spacetimi/timi_shared_server/v2/code/core/controller"
)

type VersionHandler struct { // Implements IRouteHandler
}

func (vh *VersionHandler) Routes() []controller.Route {
	return []controller.Route{
		controller.NewRoute(shared_routes.Version, []controller.RequestMethodType{controller.GET, controller.POST}),
	}
}

func (vh *VersionHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {
	_, _ = fmt.Fprintln(httpResponseWriter, "TODO")
}
