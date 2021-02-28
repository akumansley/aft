package handlers

import (
	"net/http"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
)

type APIHandler struct {
	DB db.DB
}

type apiMethod struct {
	needsRWTx          bool
	useSummaryResponse bool
}

var apiMethods = map[string]apiMethod{
	"findMany":   apiMethod{needsRWTx: false, useSummaryResponse: false},
	"findOne":    apiMethod{needsRWTx: false, useSummaryResponse: false},
	"count":      apiMethod{needsRWTx: false, useSummaryResponse: true},
	"create":     apiMethod{needsRWTx: true, useSummaryResponse: false},
	"update":     apiMethod{needsRWTx: true, useSummaryResponse: false},
	"updateMany": apiMethod{needsRWTx: true, useSummaryResponse: true},
	"delete":     apiMethod{needsRWTx: true, useSummaryResponse: false},
	"deleteMany": apiMethod{needsRWTx: true, useSummaryResponse: true},
	"upsert":     apiMethod{needsRWTx: true, useSummaryResponse: false},
}

var responseType = map[string]bool{}

func (a APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	methodName, modelName, body, err := unpackArgs(r)
	if err != nil {
		return err
	}
	m := apiMethods[methodName]

	var tx db.Tx
	if m.needsRWTx {
		tx = a.DB.NewRWTxWithContext(r.Context())
	} else {
		tx = a.DB.NewTxWithContext(r.Context())
	}

	out, err := auth.AuthedCall(tx, methodName, []interface{}{modelName, body})

	if err != nil {
		if m.needsRWTx {
			tx.Abort(err)
		}
		return err
	}

	if m.needsRWTx {
		tx.Commit()
	}

	if m.useSummaryResponse {
		response(w, &SummaryResponse{Count: out.(int)})
	} else {
		response(w, &DataResponse{Data: out})
	}
	return
}
