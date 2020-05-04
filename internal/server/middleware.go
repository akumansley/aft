package server

import (
	"awans.org/aft/internal/oplog"
	"context"
	"github.com/json-iterator/go"
	"log"
	"net/http"
	"time"
)

func Middleware(op Operation, log oplog.OpLog) http.Handler {
	// this goes inside out
	// invoke the server
	server := op.Server

	// and audit log it
	auditLoggedServer := AuditLog(op, server, log)

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
	log   oplog.OpLog
}

var apiRequestKey = "ApiRequestId"

func NewContext(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, apiRequestKey, id)
}

func ApiRequestId(ctx context.Context) uint {
	iv := ctx.Value(apiRequestKey)
	id, ok := iv.(uint)
	if !ok {
		panic("No apirequestid in context")
	}
	return id
}

// just a way to pass the id from parse->serve
type auditLoggedRequest struct {
	ApiRequestId uint
	inner        interface{}
}

func (a AuditLoggedServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	apiRequestId := a.log.NextId()
	ctx = NewContext(ctx, apiRequestId)

	pr, err := a.inner.Parse(ctx, req)
	return auditLoggedRequest{inner: pr, ApiRequestId: apiRequestId}, err
}

func (a AuditLoggedServer) Serve(ctx context.Context, req interface{}) (resp interface{}, err error) {
	al, ok := req.(auditLoggedRequest)
	if !ok {
		panic("some middleware messing with audit logger?")
	}
	ctx = NewContext(ctx, al.ApiRequestId)
	resp, err = a.inner.Serve(ctx, al.inner)
	if err == nil {
		a.log.Log(req)
	}
	return
}

func AuditLog(op Operation, inner Server, log oplog.OpLog) Server {
	return AuditLoggedServer{inner: inner, log: log}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ServerToHandler(server Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var status int
		var bytes []byte
		parsed, err := server.Parse(r.Context(), r)
		if err != nil {
			er := ErrorResponse{
				Code:    "parse-error",
				Message: err.Error(),
			}
			bytes, _ = jsoniter.Marshal(&er)
			status = http.StatusBadRequest
		} else {
			resp, err := server.Serve(r.Context(), parsed)
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
