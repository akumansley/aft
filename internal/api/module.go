package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
	"net/http"
)

type Module struct {
	lib.BlankModule
	db db.DB
	b  *bus.EventBus
}

func (m Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "FindMany",
			Pattern: "/api/{modelName}.findMany",
			Handler: errorHandler(FindManyHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "FindOne",
			Pattern: "/api/{modelName}.findOne",
			Handler: errorHandler(FindOneHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "Create",
			Pattern: "/api/{modelName}.create",
			Handler: errorHandler(CreateHandler{db: m.db, bus: m.b}),
		},
	}
}

func GetModule(db db.DB, b *bus.EventBus) lib.Module {
	return Module{db: db, b: b}
}

type apiHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorHandler(inner apiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := inner.ServeHTTP(w, r)
		if err != nil {
			er := ErrorResponse{
				Code:    "serve-error",
				Message: err.Error(),
			}
			bytes, _ := jsoniter.Marshal(&er)
			status := http.StatusBadRequest

			_, _ = w.Write(bytes)
			w.WriteHeader(status)
		}
	})
}
