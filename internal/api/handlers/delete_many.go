package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type DeleteManyHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s DeleteManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, dmBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	rwtx := s.db.NewRWTx()
	ctx := db.WithRWTx(r.Context(), rwtx)

	out, err := functions.DeleteMany([]interface{}{ctx, modelName, dmBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &SummaryResponse{Count: out.(int)})
	return
}
