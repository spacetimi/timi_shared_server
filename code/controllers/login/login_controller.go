package login

import (
	"github.com/spacetimi/timi_shared_server/code/core/controller"
)

type LoginController struct { // Implements IAppController
}

func (lc *LoginController) RouteHandlers() []controller.IRouteHandler {
    return []controller.IRouteHandler {
        &LoginHandler{},
    }
}
