package admin_server

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/controllers/admin"
    "log"
    "net/http"
)

func StartServer() {
    router := mux.NewRouter()

    router.HandleFunc("/admin", admin.AdminController).Methods("GET", "POST")
    router.HandleFunc("/admin/", admin.AdminController).Methods("GET", "POST")
    router.HandleFunc("/admin/{param1}", admin.AdminController).Methods("GET", "POST")

    // Set up static file-server for images
    router.PathPrefix("/images/").
        Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(config.GetImageFilesPath()))))

    fmt.Println("Admin Server Started on port 8001")
    log.Fatal(http.ListenAndServe(":8001", router))
}
