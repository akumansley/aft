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

	// First extract the records
	var where map[string]interface{}
	var ok bool
	if where, ok = urBody["where"].(map[string]interface{}); ok {
		delete(urBody, "where")
	}
	fmr := make(map[string]interface{})
	fmr["where"] = where

	fi, err := p.ParseFindMany(modelName, fmr)
	if err != nil {
		return
	}

	qrs, err := fi.Apply(tx)
	if err != nil {
		return
	}
	var recs []db.Record
	for _, qr := range qrs {
		recs = append(recs, qr.Record)
	}

	//Now update them
	op, err := p.ParseUpdateMany(recs, urBody)
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
