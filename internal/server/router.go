package server

import (
	"awans.org/aft/internal/server/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	s := router.Methods("POST").Subrouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Middleware(handler, route.Name)
		path := fmt.Sprintf("/api/%s.%s", route.Service, route.Method)
		s.Path(path).Name(route.Name).Handler(handler)
	}
	return router
}
