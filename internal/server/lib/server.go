package lib

import (
	"context"
	"net/http"
)

type Server interface {
	Parse(context.Context, *http.Request) (interface{}, error)
	Serve(context.Context, interface{}) (interface{}, error)
}
