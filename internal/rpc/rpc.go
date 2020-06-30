package rpc

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type RPCRequest struct {
	Data map[string]interface{} `json:"data"`
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
	RPCOut, err := eval(name, rr.Data, rwtx)
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
