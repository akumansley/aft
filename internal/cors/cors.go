package cors

import (
	"errors"
	"net/http"

	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{CSRF, CORS}
}

func CSRF(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cType := r.Header.Get("Content-Type")
		if cType != "application/json" && r.Method != "GET" {
			lib.WriteError(w, errors.New("Invalid Content-Type; must be application/json"))
			return
		}
		inner.ServeHTTP(w, r)
	})
}

// TODO check some list of acceptable origins
func CORS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}
		inner.ServeHTTP(w, r)
	})
}

func GetModule() lib.Module {
	return &Module{}
}
