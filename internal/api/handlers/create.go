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

	rwtx := s.DB.NewRWTx()
	ctx := db.WithRWTx(r.Context(), rwtx)
	out, err := functions.Create([]interface{}{ctx, modelName, crBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
