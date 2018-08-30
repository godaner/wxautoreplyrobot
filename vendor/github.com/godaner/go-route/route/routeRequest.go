package route

import "net/http"

type RouteRequest struct {
	Request *http.Request
	Params map[string]interface{}
}


