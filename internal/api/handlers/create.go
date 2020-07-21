package handlers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/api/parsers"
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
	Operation operations.CreateOperation
	Include   operations.Include
}

type CreateResponse struct {
	Data interface{} `json:"data"`
}

type CreateHandler struct {
	DB  db.DB
	Bus *bus.EventBus
}

func (s CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.DB.NewRWTx()
	p := parsers.Parser{Tx: tx}
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
	s.Bus.Publish(lib.ParseRequest{Request: request})

	st, err := request.Operation.Apply(tx)
	if err != nil {
		return
	}
	responseData := request.Include.ResolveOne(tx, st.Interface().ID(), st)
	response := CreateResponse{Data: responseData}
	tx.Commit()

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
