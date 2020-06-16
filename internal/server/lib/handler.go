package lib

import (
	"github.com/json-iterator/go"
	"net/http"
)

type ApiHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(inner ApiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := inner.ServeHTTP(w, r)
		if err != nil {
			er := ErrorResponse{
				Code:    "serve-error",
				Message: err.Error(),
			}
			bytes, _ := jsoniter.Marshal(&er)
			status := http.StatusBadRequest

			_, _ = w.Write(bytes)
			w.WriteHeader(status)
		}
	})
}
