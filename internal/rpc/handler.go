package rpc

import (
	"encoding/gob"
	"io/ioutil"
	"net/http"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type RPCResponse struct {
	Data interface{} `json:"data"`
}

type RPCHandler struct {
	bus    *bus.EventBus
	db     db.DB
	authed bool
}

func init() {
	gob.Register(map[string]interface{}{})
}

func (rh RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	name := vars["name"]
	var args map[string]interface{}
	buf, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(buf, &args)
	if err != nil {
		return
	}

	rh.bus.Publish(lib.ParseRequest{Request: args})

	ctx := r.Context()
	rwtx := rh.db.NewRWTxWithContext(ctx)

	// TODO this is silly; rewrite tx to accept context
	ctx = db.WithRWTx(ctx, rwtx)
	rwtx.SetContext(ctx)

	var RPCOut interface{}
	var f db.Function

	if rh.authed {
		RPCOut, err = auth.AuthedCall(rwtx, name, []interface{}{args})
	} else {
		function := rwtx.Ref(db.FunctionInterface.ID())
		var frec db.Record
		frec, err = rwtx.Query(function,
			db.Filter(function, db.Eq("name", name)),
			db.Filter(function, db.Eq("funcType", uuid.UUID(db.RPC.ID()))),
		).OneRecord()
		if err != nil {
			return
		}
		f, err = rwtx.Schema().LoadFunction(frec)
		if err != nil {
			return
		}
		RPCOut, err = f.Call(ctx, []interface{}{args})
	}
	if err != nil {
		return
	}
	rwtx.Commit()
	response := RPCResponse{Data: RPCOut}

	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
