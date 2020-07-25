package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"net/http"
)

type DeleteManyHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s DeleteManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, drBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewRWTx()
	p := parsers.Parser{Tx: tx}

	op, err := p.ParseDeleteMany(modelName, drBody)
	if err != nil {
		return
	}

	s.bus.Publish(lib.ParseRequest{Request: op})

	out, err := op.Apply(tx)
	if err != nil {
		return
	}
	tx.Commit()

	bytes, _ := jsoniter.Marshal(&SummaryResponse{BatchPayload{Count: out}})
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
