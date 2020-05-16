package middleware

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"context"
	"net/http"
)

type txServer struct {
	inner lib.Server
	db    db.DB
}

type rwtxServer struct {
	inner lib.Server
	db    db.DB
}

var txKey = "Tx"
var rwtxKey = "RWTx"

func NewTxContext(ctx context.Context, tx db.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func NewRWTxContext(ctx context.Context, rwtx db.RWTx) context.Context {
	return context.WithValue(ctx, rwtxKey, rwtx)
}

func TxFromContext(ctx context.Context) db.Tx {
	iv := ctx.Value(txKey)
	tx, ok := iv.(db.Tx)
	if !ok {
		panic("No tx in context")
	}
	return tx
}

func RWTxFromContext(ctx context.Context) db.RWTx {
	iv := ctx.Value(rwtxKey)
	tx, ok := iv.(db.RWTx)
	if !ok {
		panic("No tx in context")
	}
	return tx
}

// just a way to pass the tx from parse->serve
type txRequest struct {
	tx    db.Tx
	inner interface{}
}
type rwtxRequest struct {
	rwtx  db.RWTx
	inner interface{}
}

func (t txServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	tx := t.db.NewTx()
	ctx = NewTxContext(ctx, tx)
	pr, err := t.inner.Parse(ctx, req)
	return txRequest{inner: pr, tx: tx}, err
}

func (t txServer) Serve(ctx context.Context, req interface{}) (resp interface{}, err error) {
	txr, ok := req.(txRequest)
	if !ok {
		panic("some middleware messing with tx ?")
	}
	tx := txr.tx.(db.Tx)
	ctx = NewTxContext(ctx, tx)
	resp, err = t.inner.Serve(ctx, txr.inner)
	return
}

func (t rwtxServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	rwtx := t.db.NewRWTx()
	ctx = NewRWTxContext(ctx, rwtx)
	pr, err := t.inner.Parse(ctx, req)
	return rwtxRequest{inner: pr, rwtx: rwtx}, err
}

func (t rwtxServer) Serve(ctx context.Context, req interface{}) (resp interface{}, err error) {
	txr, ok := req.(rwtxRequest)
	rwtx := txr.rwtx.(db.RWTx)
	ctx = NewRWTxContext(ctx, rwtx)
	if !ok {
		panic("some middleware messing with tx ?")
	}
	resp, err = t.inner.Serve(ctx, txr.inner)
	if err == nil {
		err = rwtx.Commit()
	}
	return
}

func Tx(db db.DB, inner lib.Server) lib.Server {
	return txServer{inner: inner, db: db}
}

func RWTx(db db.DB, inner lib.Server) lib.Server {
	return rwtxServer{inner: inner, db: db}
}
