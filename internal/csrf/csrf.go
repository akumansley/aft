package csrf

import (
	"errors"
	"net/http"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("738f74a2-6d1f-4633-9f60-b3f47914cd2c")
}

func (m *Module) Name() string {
	return "csrf"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/csrf"
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{CSRF}
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

func GetModule() lib.Module {
	return &Module{}
}
