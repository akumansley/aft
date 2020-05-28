package cors

import (
	"awans.org/aft/internal/server/lib"
	"net/http"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{CORS}
}

func CORS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		inner.ServeHTTP(w, r)
	})
}

func GetModule() lib.Module {
	return &Module{}
}
