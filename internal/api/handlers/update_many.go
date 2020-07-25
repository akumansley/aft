package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
	"net/http"
)

type UpdateManyHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpdateManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, urBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewRWTx()
	p := parsers.Parser{Tx: tx}

	op, err := p.ParseUpdateMany(modelName, urBody)
	if err != nil {
		return
	}

	s.bus.Publish(lib.ParseRequest{Request: op})

	out, err := op.Apply(tx)
	if err != nil {
		return
	}
	tx.Commit()

	bytes, _ := jsoniter.Marshal(&SummaryResponse{Count: out})
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
