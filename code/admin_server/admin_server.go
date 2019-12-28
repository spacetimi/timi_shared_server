package admin_server

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/spacetimi/timi_shared_server/code/controllers/admin"
    "log"
    "net/http"
)

func StartServer() {
    router := mux.NewRouter()

    router.HandleFunc("/admin", admin.AdminController).Methods("GET", "POST")
    router.HandleFunc("/admin/", admin.AdminController).Methods("GET", "POST")
    router.HandleFunc("/admin/{param1}", admin.AdminController).Methods("GET", "POST")

    fmt.Println("Admin Server Started on port 9000")
    log.Fatal(http.ListenAndServe(":9000", router))
}
