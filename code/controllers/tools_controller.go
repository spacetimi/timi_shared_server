package controllers

import (
	"fmt"
	"net/http"
)

func ToolsController(httpResponseWriter http.ResponseWriter, request *http.Request) {
	// TODO: Route to different tools, preferably using angular js
	fmt.Fprintf(httpResponseWriter, "Tools")
}
