package server

import (
	"log"
	"net/http"
	"time"
)

func Middleware(op Operation) http.Handler {
	// this goes inside out
	// invoke the server
	server := op.Server

	// but oplog it
	opLoggedServer := OpLog(op, server)

	// and audit log it
	auditLoggedServer := AuditLog(op, opLoggedServer)

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

func (a AuditLoggedServer) Parse(req *http.Request) interface{} {
	return a.inner.Parse(req)
}

func (a AuditLoggedServer) Serve(w http.ResponseWriter, req interface{}) {
	a.inner.Serve(w, req)
}

func AuditLog(op Operation, inner Server) Server {
	return AuditLoggedServer{inner: inner}
}

type OpLoggedServer struct {
	inner Server
}

func (o OpLoggedServer) Parse(req *http.Request) interface{} {
	return o.inner.Parse(req)
}

func (o OpLoggedServer) Serve(w http.ResponseWriter, req interface{}) {
	o.inner.Serve(w, req)
}

func OpLog(op Operation, inner Server) Server {
	if op.Write {
		return OpLoggedServer{inner: inner}
	}
	return inner
}

func ServerToHandler(server Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsed := server.Parse(r)
		server.Serve(w, parsed)
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
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
