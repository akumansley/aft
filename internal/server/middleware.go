package server

import (
	"github.com/json-iterator/go"
	"log"
	"net/http"
	"time"
)

func Middleware(op Operation) http.Handler {
	// this goes inside out
	// invoke the server
	server := op.Server

	// and audit log it
	auditLoggedServer := AuditLog(op, server)

	// wrap it as a handler
	handler := ServerToHandler(auditLoggedServer)

	// before that we log the request
	logged := Logger(handler, op.Name)

	// before that we set CORS
	cors := CORS(logged)

	return cors
}

type AuditLoggedServer struct {
	inner Server
}

func (a AuditLoggedServer) Parse(req *http.Request) (interface{}, error) {
	pr, err := a.inner.Parse(req)
	bytes, _ := jsoniter.Marshal(&pr)
	log.Printf(
		"%v",
		string(bytes),
	)
	return pr, err
}

func (a AuditLoggedServer) Serve(req interface{}) (interface{}, error) {
	return a.inner.Serve(req)
}

func AuditLog(op Operation, inner Server) Server {
	return AuditLoggedServer{inner: inner}
}

type ErrorResponse struct {
	Code    string
	Message string
}

func ServerToHandler(server Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var status int
		var bytes []byte
		parsed, err := server.Parse(r)
		if err != nil {
			er := ErrorResponse{
				Code:    "parse-error",
				Message: err.Error(),
			}
			bytes, _ = jsoniter.Marshal(&er)
			status = http.StatusBadRequest
		}
		resp, err := server.Serve(parsed)
		if err != nil {
			er := ErrorResponse{
				Code:    "serve-error",
				Message: err.Error(),
			}
			bytes, _ = jsoniter.Marshal(&er)
			status = http.StatusBadRequest
		} else {
			bytes, _ = jsoniter.Marshal(&resp)
			status = http.StatusOK
		}

		_, _ = w.Write(bytes)
		w.WriteHeader(status)
	})
}

func CORS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		inner.ServeHTTP(w, r)
	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%-30.30s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
