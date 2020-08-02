package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"net/http"
)

type UpsertHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpsertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, usrBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewRWTx()
	p := parsers.Parser{Tx: tx}

	op, err := p.ParseUpsert(modelName, usrBody)
	if err != nil {
		return
	}

	s.bus.Publish(lib.ParseRequest{Request: op})

	out, err := op.Apply(tx)
	if err != nil {
		return
	}
	tx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
