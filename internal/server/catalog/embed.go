package catalog

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
)

var pkgDir = pkger.Dir("/client/catalog/public/")

type fs pkger.Dir

func (f fs) Stat(name string) (os.FileInfo, error) {
	dirPath := string(f)
	path := filepath.Join(dirPath, name)
	return pkger.Stat(path)
}

func (f fs) Open(name string) (http.File, error) {
	return pkger.Dir(f).Open(name)
}

var Dir = fs(pkgDir)
