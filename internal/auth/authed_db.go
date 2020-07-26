package auth

import (
	"awans.org/aft/internal/db"
)

type authedDB struct {
	db.DB
}

type authedRWTx struct {
	db.RWTx
}

type authedTx struct {
	db.Tx
}

func AuthedDB(d db.DB) db.DB {
	return &authedDB{DB: d}
}

func (d *authedDB) NewRWTx() db.RWTx {
	return &authedRWTx{RWTx: d.DB.NewRWTx()}
}

func (d *authedDB) NewTx() db.Tx {
	return &authedTx{Tx: d.DB.NewTx()}
}

func (t *authedTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	panic("Not implemented")
}

func (t *authedRWTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	panic("Not implemented")
}
