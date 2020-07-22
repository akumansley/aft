package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
	"net/http"
)

type CreateHandler struct {
	DB  db.DB
	Bus *bus.EventBus
}

func (s CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, crBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.DB.NewRWTx()
	p := parsers.Parser{Tx: tx}

	op, err := p.ParseCreate(modelName, crBody)
	if err != nil {
		return
	}

	s.Bus.Publish(lib.ParseRequest{Request: op})

	out, err := op.Apply(tx)
	if err != nil {
		return
	}
	tx.Commit()

	bytes, _ := jsoniter.Marshal(&DataResponse{Data: out})
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
