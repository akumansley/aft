package aft

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed client/catalog/public
var rootFS embed.FS

func init() {
	subFS, err := fs.Sub(rootFS, "client/catalog/public")
	if err != nil {
		panic(err)
	}
	Catalog = http.FS(subFS)
}

var Catalog http.FileSystem
