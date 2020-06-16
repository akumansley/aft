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

type UpdateManyRequestBody struct {
	Where map[string]interface{} `json:"where"`
	Data  map[string]interface{} `json:"data"`
}

type UpdateManyRequest struct {
	Operation UpdateManyOperation
}
type BatchPayload struct {
	Count int `json:"count"`
}
type UpdateManyResponse struct {
	BatchPayload `json:"BatchPayload"`
}

type UpdateManyHandler struct {
	DB  db.DB
	Bus *bus.EventBus
}

func (s UpdateManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.DB.NewRWTx()
	p := Parser{tx: tx}
	var urBody UpdateManyRequestBody
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	body, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(body, &urBody)
	if err != nil {
		return
	}
	//find the records to update
	var firequest FindManyRequest
	fi, err := p.ParseFindMany(modelName, urBody.Where)
	if err != nil {
		return
	}
	firequest = FindManyRequest{
		Operation: fi,
	}
	s.Bus.Publish(lib.ParseRequest{Request: firequest})

	recs := firequest.Operation.Apply(tx)

	//Now update the record
	var request UpdateManyRequest
	op, err := p.ParseUpdateMany(recs, urBody.Data)
	if err != nil {
		return
	}
	request = UpdateManyRequest{
		Operation: op,
	}
	s.Bus.Publish(lib.ParseRequest{Request: request})
	st, err := request.Operation.Apply(tx)
	if err != nil {
		return
	}

	response := UpdateManyResponse{BatchPayload{Count: st}}
	tx.Commit()

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
