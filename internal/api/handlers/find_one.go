package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
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

	bytes, _ := jsoniter.Marshal(&DataResponse{Data: out})
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
