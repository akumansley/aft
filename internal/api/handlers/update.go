package handlers

import (
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
	"net/http"
)

type UpdateHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, urBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewRWTx()
	p := parsers.Parser{Tx: tx}

	//first get the record
	var where map[string]interface{}
	var ok bool
	if where, ok = urBody["where"].(map[string]interface{}); ok {
		delete(urBody, "where")
	}

	fir := make(map[string]interface{})
	fir["where"] = where
	fi, err := p.ParseFindOne(modelName, fir)
	if err != nil {
		return
	}
	rec, err := fi.Apply(tx)
	if err != nil {
		return
	}

	//Now update it
	op, err := p.ParseUpdate(rec.Record, urBody)
	if err != nil {
		return
	}

	s.bus.Publish(lib.ParseRequest{Request: op})

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
