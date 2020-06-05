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

type CreateRequestBody struct {
	Data    map[string]interface{} `json:"data"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type CreateRequest struct {
	// TODO add Select
	Operation CreateOperation
	Include   Include
}

type CreateResponse struct {
	Data interface{} `json:"data"`
}

type CreateHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.db.NewRWTx()
	p := Parser{tx: tx}
	var crBody CreateRequestBody
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	body, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(body, &crBody)
	if err != nil {
		return
	}
	var request CreateRequest
	op, err := p.ParseCreate(modelName, crBody.Data)
	if err != nil {
		return
	}
	inc, err := p.ParseInclude(modelName, crBody.Include)
	if err != nil {
		return
	}
	request = CreateRequest{
		Operation: op,
		Include:   inc,
	}
	s.bus.Publish(lib.ParseRequest{Request: request})

	st, err := request.Operation.Apply(tx)
	if err != nil {
		return
	}
	responseData := request.Include.Resolve(tx, st)
	response := CreateResponse{Data: responseData}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
