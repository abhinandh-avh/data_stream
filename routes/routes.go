package routes

import (
	"net/http"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

// RouterImpl is an implementation of the Router interface.
type RouterImpl struct {
	routes []route
}

// NewRouter creates a new RouterImpl instance.
func NewRouter() *RouterImpl {
	return &RouterImpl{}
}

// AddRoute adds a new route to the router.
func (r *RouterImpl) AddRoute(method, path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{method, path, handler})
}

// ServeHTTP handles incoming HTTP requests by routing them to the appropriate handler.
func (r *RouterImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method == route.method && req.URL.Path == route.path {
			route.handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}
