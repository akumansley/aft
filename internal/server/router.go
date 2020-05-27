package server

import (
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	"net/http"
)

type router struct {
	r          *mux.Router
	entrypoint http.Handler
}

func NewRouter() *router {
	r := router{}
	r.r = mux.NewRouter().Methods("POST").Subrouter()
	r.entrypoint = r.r
	return &r
}

func (r *router) AddMiddleware(middleware []lib.Middleware) {
	for _, m := range middleware {
		r.entrypoint = m(r.entrypoint)
	}
}

func (r *router) AddRoutes(routes []lib.Route) {
	for _, route := range routes {
		r.r.Path(route.Pattern).Name(route.Name).Handler(route.Handler)
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.entrypoint.ServeHTTP(w, req)
}
