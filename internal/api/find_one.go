package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type FindOneRequestBody struct {
	Where   map[string]interface{} `json:"where"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type FindOneRequest struct {
	// TODO add Select
	Operation FindOneOperation
	Include   Include
}

type FindOneResponse struct {
	Data interface{} `json:"data"`
}

type FindOneHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s FindOneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.db.NewTx()
	p := Parser{tx: tx}
	var foBody FindOneRequestBody
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	buf, _ := ioutil.ReadAll(r.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)
	op, err := p.ParseFindOne(modelName, foBody.Where)
	if err != nil {
		return
	}
	inc, err := p.ParseInclude(modelName, foBody.Include)
	if err != nil {
		return
	}

	request := FindOneRequest{
		Operation: op,
		Include:   inc,
	}

	s.bus.Publish(lib.ParseRequest{Request: request})

	st, err := request.Operation.Apply(tx)
	if err != nil {
		return
	}
	responseData := request.Include.ResolveOne(tx, op.ModelID, st)
	response := FindOneResponse{Data: responseData}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
