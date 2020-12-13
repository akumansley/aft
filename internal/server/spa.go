package server

import (
	"net/http"
	"os"
	"path/filepath"
)

// modified from https://github.com/gorilla/mux/README.md

type OpenStater interface {
	Open(name string) (http.File, error)
	Stat(name string) (os.FileInfo, error)
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	Dir OpenStater
}

const indexPath = "index.html"

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check whether a file exists at the given path
	_, err = h.Dir.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		f, err := h.Dir.Open(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		finfo, err := f.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.ServeContent(w, r, indexPath, finfo.ModTime(), f)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(h.Dir).ServeHTTP(w, r)
}
