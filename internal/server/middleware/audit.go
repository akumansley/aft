package middleware

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
	"context"
	"github.com/google/uuid"
	"net/http"
)

type auditLoggedServer struct {
	inner lib.Server
	log   oplog.OpLog
}

var apiRequestKey = "ApiRequestId"

func newContext(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, apiRequestKey, id)
}

func ApiRequestIdFromContext(ctx context.Context) uuid.UUID {
	iv := ctx.Value(apiRequestKey)
	id, ok := iv.(uuid.UUID)
	if !ok {
		panic("No apirequestid in context")
	}
	return id
}

// just a way to pass the id from parse->serve
type auditLoggedRequest struct {
	ApiRequestId uuid.UUID
	inner        interface{}
}

func (a auditLoggedServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	apiRequestId := uuid.New()
	ctx = newContext(ctx, apiRequestId)

	pr, err := a.inner.Parse(ctx, req)
	return auditLoggedRequest{inner: pr, ApiRequestId: apiRequestId}, err
}

func (a auditLoggedServer) Serve(ctx context.Context, req interface{}) (resp interface{}, err error) {
	al, ok := req.(auditLoggedRequest)
	if !ok {
		panic("some middleware messing with audit logger?")
	}
	ctx = newContext(ctx, al.ApiRequestId)
	resp, err = a.inner.Serve(ctx, al.inner)
	if err == nil {
		a.log.Log(req)
	}
	return
}

func AuditLog(op lib.Operation, inner lib.Server, log oplog.OpLog) lib.Server {
	return auditLoggedServer{inner: inner, log: log}
}
