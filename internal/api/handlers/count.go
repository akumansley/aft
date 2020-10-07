package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type CountHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s CountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, fmBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	tx := s.db.NewTx()
	ctx := db.WithTx(r.Context(), tx)

	out, err := functions.Count.Call([]interface{}{ctx, modelName, fmBody})
	if err != nil {
		return err
	}

	response(w, &SummaryResponse{Count: out.(int)})
	return
}
