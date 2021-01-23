package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
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

	tx := s.db.NewTx()
	ctx := db.WithTx(r.Context(), tx)

	out, err := functions.FindOne([]interface{}{ctx, modelName, foBody})
	if err != nil {
		return err
	}

	response(w, &DataResponse{Data: out})
	return
}
