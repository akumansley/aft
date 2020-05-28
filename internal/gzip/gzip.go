package gzip

import (
	"awans.org/aft/internal/server/lib"
	"github.com/NYTimes/gziphandler"
	"net/http"
)

type Module struct {
	lib.BlankModule
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
