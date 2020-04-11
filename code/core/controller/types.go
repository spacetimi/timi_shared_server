package controller

import "net/http"

type IAppController interface {
    RouteHandlers() []IRouteHandler
}

type IRouteHandler interface {
    Routes() []Route
    HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *HandlerFuncArgs)
}

type Route struct {
    Path string
    Methods []RequestMethodType
}

func NewRoute(path string, methods []RequestMethodType) Route {
    return Route {
        Path:path,
        Methods:methods,
    }
}

func (r Route) GetMethodsAsStrings() []string {
    var s []string
    for _, method := range r.Methods {
        s = append(s, method.String())
    }
    return s
}

type HandlerFuncArgs struct {
    RequestPathVars map[string]string
    PostArgs map[string]string
}

type RequestMethodType int
const (
    GET RequestMethodType = iota
    POST
    PUT
)

func (rmt RequestMethodType)String() string {
	switch rmt {
	case GET:
		return "GET"
	case POST:
		return "POST"
    case PUT:
        return "PUT"
	}
	return "GET"
}

