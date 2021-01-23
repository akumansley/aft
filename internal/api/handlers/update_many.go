package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type UpdateManyHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpdateManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, umBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	rwtx := s.db.NewRWTx()
	ctx := db.WithRWTx(r.Context(), rwtx)

	out, err := functions.UpdateMany.Call([]interface{}{ctx, modelName, umBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &SummaryResponse{Count: out.(int)})
	return
}