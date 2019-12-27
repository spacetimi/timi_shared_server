package api_server

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/spacetimi/timi_shared_server/code/controllers"
    "github.com/spacetimi/timi_shared_server/code/controllers/login"
    "log"
    "net/http"
)


func StartServer(testingController func(w http.ResponseWriter, response *http.Request),
                 dummyController func(w http.ResponseWriter, response *http.Request)) {

    router := mux.NewRouter()

    router.HandleFunc("/login", login.HandleLogin).Methods("POST")

    router.HandleFunc("/", controllers.TestController).Methods("GET", "POST")
    router.HandleFunc("/test", controllers.TestController).Methods("GET", "POST")
    router.HandleFunc("/redis-ping", controllers.PingRedisController).Methods("GET", "POST")
    router.HandleFunc("/tools", controllers.ToolsController).Methods("GET", "POST")
    router.HandleFunc("/fatal", controllers.FatalController).Methods("GET", "POST")
    router.HandleFunc("/panic", controllers.PanicController).Methods("GET", "POST")

    router.HandleFunc("/testing/{param1}", testingController).Methods("GET", "POST")
    router.HandleFunc("/testing/{param1}/{param2}", testingController).Methods("GET", "POST")
    router.HandleFunc("/testing/{param1}/{param2}/{param3}", testingController).Methods("GET", "POST")

    router.HandleFunc("/dummy/{param1}", dummyController).Methods("GET", "POST")
    router.HandleFunc("/testing/{param1}/{param2}", testingController).Methods("GET", "POST")
    router.HandleFunc("/testing/{param1}/{param2}/{param3}", testingController).Methods("GET", "POST")


    fmt.Println("Server started on port 8000")
    log.Fatal(http.ListenAndServe(":8000", router))

}

