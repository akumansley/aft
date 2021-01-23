package rpc

import (
	"encoding/gob"
	"io/ioutil"
	"net/http"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type RPCResponse struct {
	Data interface{} `json:"data"`
}

type RPCHandler struct {
	bus *bus.EventBus
	db  db.DB
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

	rwtx := rh.db.NewRWTxWithContext(r.Context())
	ctx := db.WithRWTx(r.Context(), rwtx)

	f, role, err := loadFunction(rwtx, name)
	if err != nil {
		return
	}

	if role != nil {
		// restart the tx with the new role
		// is this cool
		rwtx.Commit()
		ctx = auth.WithRole(r.Context(), role)
		rwtx = rh.db.NewRWTxWithContext(ctx)
		ctx = db.WithRWTx(ctx, rwtx)
	}

	RPCOut, err := f.Call([]interface{}{ctx, args})
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

func loadFunction(tx db.RWTx, name string) (f db.Function, role db.Record, err error) {
	rpcs := tx.Ref(RPCModel.ID())
	function := tx.Ref(db.FunctionInterface.ID())
	roles := tx.Ref(auth.RoleModel.ID())
	rpcFunction, _ := tx.Schema().GetRelationshipByID(RPCFunction.ID())
	rpcRole, _ := tx.Schema().GetRelationshipByID(RPCRole.ID())

	result, err := tx.Query(rpcs,
		db.Join(function, rpcs.Rel(rpcFunction)),
		db.Filter(function, db.Eq("name", name)),
		db.LeftJoin(roles, rpcs.Rel(rpcRole)),
	).One()
	if err != nil {
		return
	}
	funcRec := result.GetChildRelOne(rpcFunction).Record
	roleQR := result.GetChildRelOne(rpcRole)
	if roleQR != nil {
		role = roleQR.Record
	}

	f, err = tx.Schema().LoadFunction(funcRec)
	return
}
