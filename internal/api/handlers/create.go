package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
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

	rwtx := s.DB.NewRWTxWithContext(r.Context())
	ctx := db.WithRWTx(r.Context(), rwtx)
	out, err := functions.Create(ctx, []interface{}{modelName, crBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
