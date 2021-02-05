package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type DeleteHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, dlBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	rwtx := s.db.NewRWTxWithContext(r.Context())
	ctx := db.WithRWTx(r.Context(), rwtx)

	out, err := functions.Delete(ctx, []interface{}{modelName, dlBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
