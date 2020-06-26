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

type UpdateRequestBody struct {
	Where   map[string]interface{} `json:"where"`
	Data    map[string]interface{} `json:"data"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type UpdateRequest struct {
	// TODO add Select
	Operation UpdateOperation
	Include   Include
}

type UpdateResponse struct {
	Data interface{} `json:"data"`
}

type UpdateHandler struct {
	DB  db.DB
	Bus *bus.EventBus
}

func (s UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.DB.NewRWTx()
	p := Parser{tx: tx}
	var urBody UpdateRequestBody
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	body, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(body, &urBody)
	if err != nil {
		return
	}
	//find the record to update
	var firequest FindOneRequest
	fi, err := p.ParseFindOne(modelName, urBody.Where)
	if err != nil {
		return
	}
	firequest = FindOneRequest{
		Operation: fi,
	}
	s.Bus.Publish(lib.ParseRequest{Request: firequest})

	rec, err := firequest.Operation.Apply(tx)
	if err != nil {
		return
	}

	//Now update the record
	var request UpdateRequest
	op, err := p.ParseUpdate(rec, urBody.Data)
	if err != nil {
		return
	}
	inc, err := p.ParseInclude(modelName, urBody.Include)
	if err != nil {
		return
	}
	request = UpdateRequest{
		Operation: op,
		Include:   inc,
	}
	s.Bus.Publish(lib.ParseRequest{Request: request})
	st, err := request.Operation.Apply(tx)
	if err != nil {
		return
	}

	responseData := request.Include.ResolveOne(tx, st.Model().ID, st)
	response := UpdateResponse{Data: responseData}
	tx.Commit()

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
