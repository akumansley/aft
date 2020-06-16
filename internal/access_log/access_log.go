package access_log

import (
	"awans.org/aft/internal/server/lib"
	"log"
	"net/http"
	"time"
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

		log.Printf(
			"%s\t%-30.30s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
