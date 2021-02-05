package access_log

import (
	"log"
	"net/http"
	"time"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("2e002d0c-764a-4991-be76-eded6dab1b42")
}

func (m *Module) Name() string {
	return "accessLog"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/access_log"
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

		rec, ok := auth.UserFromContext(r.Context())
		ids := "unauthed"
		if ok {
			ids = rec.ID().String()
		}

		log.Printf(
			"%-4.4s\t%-30.30s\t%-8.8s\t%s",
			r.Method,
			r.RequestURI,
			ids,
			time.Since(start),
		)
	})
}
