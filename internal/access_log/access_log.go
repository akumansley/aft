package access_log

import (
	"log"
	"net/http"
	"time"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{Logger}
}

func GetModule() lib.Module {
	return &Module{}
}

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		id, ok := auth.IDFromContext(r.Context())
		ids := "unauthed"
		if ok {
			ids = id.String()
		}

		log.Printf(
			"%s\t%-30.30s\t%-8.8s\t%s",
			r.Method,
			r.RequestURI,
			ids,
			time.Since(start),
		)
	})
}
