package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"net/http"
)

type FindOneHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s FindOneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, foBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewRWTx()
	p := parsers.Parser{Tx: tx}

	op, err := p.ParseFindOne(modelName, foBody)
	if err != nil {
		return
	}

	s.bus.Publish(lib.ParseRequest{Request: op})

	out, err := op.Apply(tx)
	if err != nil {
		return
	}

	response(w, &DataResponse{Data: out})
	return
}
