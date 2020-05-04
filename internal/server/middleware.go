package server

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/middleware"
	"context"
	"github.com/json-iterator/go"
	"log"
	"net/http"
	"time"
)

func Middleware(op Operation, db db.DB, log oplog.OpLog) http.Handler {
	// this goes inside out
	// invoke the server
	server := op.Server

	if op.Tx == RWTx {
		server = middleware.RWTx(db, server)
	} else if op.Tx == Tx {
		server = middleware.Tx(db, server)
	}

	// and audit log it
	auditLoggedServer := middleware.AuditLog(op, server, log)

	// wrap it as a handler
	handler := ServerToHandler(auditLoggedServer)

	// before that we log the request
	logged := Logger(handler, op.Name)

	// before that we set CORS
	cors := CORS(logged)

	return cors
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
