package rpc

import (
	"context"
	"io/ioutil"
	"net/http"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type RPCRequest struct {
	Args map[string]interface{} `json:"args"`
}

type RPCResponse struct {
	Data interface{} `json:"data"`
}

type RPCHandler struct {
	bus *bus.EventBus
	db  db.DB
}

func (rh RPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	name := vars["name"]
	var rr RPCRequest
	buf, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(buf, &rr)
	if err != nil {
		return
	}

	rh.bus.Publish(lib.ParseRequest{Request: rr})
	rwtx := rh.db.NewRWTx()
	ctx := db.WithRWTx(r.Context(), rwtx)

	RPCOut, err := eval(ctx, name, rr.Args, rwtx)
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

func eval(ctx context.Context, name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	rpcs := tx.Ref(RPCModel.ID())
	function := tx.Ref(db.FunctionInterface.ID())

	result, err := tx.Query(rpcs, db.Join(function, rpcs.Rel(RPCFunction)), db.Filter(function, db.Eq("name", name))).One()
	if err != nil {
		return nil, err
	}
	funcRec := result.GetChildRelOne(RPCFunction).Record

	f, err := tx.Schema().LoadFunction(funcRec)
	if err != nil {
		return nil, err
	}
	return f.Call([]interface{}{ctx, args})
}
