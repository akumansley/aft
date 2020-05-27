package lib

import "net/http"

type Route struct {
	Pattern string
	Name    string
	Handler http.Handler
}
