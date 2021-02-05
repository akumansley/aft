package gzip

import (
	"net/http"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/NYTimes/gziphandler"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("bd597db3-55db-4b38-b953-9cbca2433a25")
}

func (m *Module) Name() string {
	return "gzip"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/gzip"
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{GZipMiddleware}
}

func GZipMiddleware(inner http.Handler) http.Handler {
	return gziphandler.GzipHandler(inner)
}

func GetModule() lib.Module {
	return &Module{}
}
