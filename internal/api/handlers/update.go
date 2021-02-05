package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type UpdateHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, upBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	rwtx := s.db.NewRWTxWithContext(r.Context())
	ctx := db.WithRWTx(r.Context(), rwtx)

	out, err := functions.Update(ctx, []interface{}{modelName, upBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
