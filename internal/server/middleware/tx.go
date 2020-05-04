package middleware

import (
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/lib"
	"context"
	"net/http"
)

type txServer struct {
	inner lib.Server
	db    db.DB
	rw    bool
}

var txKey = "Tx"

func NewTxContext(ctx context.Context, tx db.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func NewRWTxContext(ctx context.Context, tx db.RWTx) context.Context {
	return context.WithValue(ctx, txKey, tx)
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
	iv := ctx.Value(txKey)
	tx, ok := iv.(db.RWTx)
	if !ok {
		panic("No tx in context")
	}
	return tx
}

// just a way to pass the tx from parse->serve
type txRequest struct {
	tx    interface{}
	inner interface{}
}

func (t txServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	var tx interface{}
	if t.rw {
		rwtx := t.db.NewRWTx()
		ctx = NewRWTxContext(ctx, rwtx)
		tx = rwtx
	} else {
		rtx := t.db.NewTx()
		ctx = NewTxContext(ctx, rtx)
		tx = rtx
	}

	pr, err := t.inner.Parse(ctx, req)
	return txRequest{inner: pr, tx: tx}, err
}

func (t txServer) Serve(ctx context.Context, req interface{}) (resp interface{}, err error) {
	txr, ok := req.(txRequest)
	if !ok {
		panic("some middleware messing with tx ?")
	}
	if t.rw {
		rwtx := txr.tx.(db.RWTx)
		ctx = NewRWTxContext(ctx, rwtx)
	} else {
		tx := txr.tx.(db.Tx)
		ctx = NewTxContext(ctx, tx)
	}
	resp, err = t.inner.Serve(ctx, txr.inner)
	if err == nil && t.rw {
		rwtx := RWTxFromContext(ctx)
		rwtx.Commit()
	}
	return
}

func Tx(db db.DB, inner lib.Server) lib.Server {
	return txServer{inner: inner, db: db, rw: false}
}

func RWTx(db db.DB, inner lib.Server) lib.Server {
	return txServer{inner: inner, db: db, rw: true}
}
