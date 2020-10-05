package cors

import (
	"net/http"

	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{CORS}
}

func CORS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
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
