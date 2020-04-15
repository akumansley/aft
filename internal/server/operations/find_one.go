package operations

import (
	"awans.org/aft/internal/server/db"
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
	// TODO add Include/Select
	Operation db.FindOneOperation
}

type FindOneResponse struct {
	Data interface{} `json:"data"`
}

type FindOneServer struct {
	DB db.DB
}

func (s FindOneServer) Parse(req *http.Request) interface{} {
	p := Parser{db: s.DB}
	var foBody FindOneRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)
	op := p.ParseFindOne(modelName, foBody.Where)

	request := FindOneRequest{
		Operation: op,
	}

	return request
}

func (s FindOneServer) Serve(w http.ResponseWriter, req interface{}) {
	params := req.(FindOneRequest)
	st := params.Operation.Apply(s.DB)
	response := FindOneResponse{Data: st}
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
}
