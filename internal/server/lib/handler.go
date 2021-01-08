package lib

import (
	"net/http"

	"awans.org/aft/internal/errors"
	jsoniter "github.com/json-iterator/go"
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
			WriteError(w, err)
		}
	})
}

func WriteError(w http.ResponseWriter, err error) {
	var code string
	if aerr, ok := err.(errors.AftError); ok {
		code = aerr.Code
	} else {
		code = "serve-error"
	}
	er := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}
	bytes, _ := jsoniter.Marshal(&er)
	status := http.StatusBadRequest

	_, _ = w.Write(bytes)
	w.WriteHeader(status)
}
