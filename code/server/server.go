package server

import (
    "github.com/gorilla/mux"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/controllers/admin"
    "github.com/spacetimi/timi_shared_server/code/controllers/login"
    "github.com/spacetimi/timi_shared_server/code/controllers/server_status"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "log"
    "net/http"
    "strconv"
)

var _router *mux.Router

/** Package init **/
func init() {
    _router = mux.NewRouter()
}

func StartServer(appController controller.IAppController) {

    if appController == nil {
        logger.LogFatal("no app controller specified")
        return
    }

    registerController(appController)
    registerController(&server_status.ServerStatusController{})
    registerController(&login.LoginController{})

    // Admin server
    _router.PathPrefix("/admin").HandlerFunc(admin.AdminController).Methods("GET", "POST")

    // Set up static file-server for images
    _router.PathPrefix("/images/").
        Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(config.GetSharedImageFilesPath()))))
    _router.PathPrefix("/app-images/").
        Handler(http.StripPrefix("/app-images/", http.FileServer(http.Dir(config.GetAppImageFilesPath()))))

    portNumberString := strconv.Itoa(config.GetEnvironmentConfiguration().Port)
    logger.LogInfo("Server started successfully|port=" + portNumberString)
    log.Fatal(http.ListenAndServe(":" + portNumberString, _router))
}

func registerController(c controller.IAppController) {
    for _, routeHandler := range c.RouteHandlers() {
        routeHandler := routeHandler
        for _, route := range routeHandler.Routes() {
            methods := route.GetMethodsAsStrings()
            _router.HandleFunc(route.Path, func(httpResponseWriter http.ResponseWriter, request *http.Request) {

                requestPathVars := mux.Vars(request)

                postArgs := make(map[string]string)
                if request.Method == controller.POST.String() {
                    err := request.ParseForm()
                    if err == nil {
                        for key, _ := range request.Form {
                            postArgs[key] = request.Form.Get(key)
                        }
                    } else {
                        logger.LogError("error parsing post args" +
                                        "|request url=" + request.URL.Path +
                                        "|error=" + err.Error())
                    }
                }

                args := &controller.HandlerFuncArgs {
                   RequestPathVars: requestPathVars,
                   PostArgs: postArgs,
                }

                routeHandler.HandlerFunc(httpResponseWriter, request, args)
            }).Methods(methods...)
        }
    }
}

