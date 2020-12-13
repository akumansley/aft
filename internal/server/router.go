package server

import (
	"net/http"

	"awans.org/aft/internal/server/catalog"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
)

type router struct {
	postRouter *mux.Router
	router     *mux.Router
	entrypoint http.Handler
}

func NewRouter() *router {
	r := router{}
	r.router = mux.NewRouter()
	r.postRouter = r.router.Methods("POST").Subrouter()
	r.entrypoint = r.router
	r.router.Methods("GET").PathPrefix("/").Handler(spaHandler{Dir: catalog.Dir})
	return &r
}

func (r *router) AddMiddleware(middleware []lib.Middleware) {
	for _, m := range middleware {
		r.entrypoint = m(r.entrypoint)
	}
}

func (r *router) AddRoutes(routes []lib.Route) {
	for _, route := range routes {
		r.postRouter.Path(route.Pattern).Name(route.Name).Handler(route.Handler)
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.entrypoint.ServeHTTP(w, req)
}
