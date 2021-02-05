package handlers

import (
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

type UpsertHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s UpsertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	modelName, usrBody, err := unpackArgs(r)
	if err != nil {
		return err
	}

	rwtx := s.db.NewRWTxWithContext(r.Context())
	ctx := db.WithRWTx(r.Context(), rwtx)

	out, err := functions.Upsert(ctx, []interface{}{modelName, usrBody})
	if err != nil {
		return err
	}

	rwtx.Commit()

	response(w, &DataResponse{Data: out})
	return
}
