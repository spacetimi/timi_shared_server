package server_status

import (
    "github.com/spacetimi/timi_shared_server/code/core/controller"
)

type ServerStatusController struct { // Implements IAppController
}

func (ssc *ServerStatusController) RouteHandlers() []controller.IRouteHandler {
    return []controller.IRouteHandler {
        &HealthCheckHandler{},
        &VersionHandler{},
    }
}
