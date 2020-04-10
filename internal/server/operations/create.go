package operations

import (
	"awans.org/aft/internal/server/db"
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
	// TODO add Include/Select
	Operation CreateOperation
}

type CreateResponse struct {
	Data interface{} `json:"data"`
}

type CreateServer struct {
	DB db.DB
}

func (s CreateServer) Parse(req *http.Request) interface{} {
	p := Parser{db: s.DB}
	var crBody CreateRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	body, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(body, &crBody)
	var request CreateRequest
	op := p.ParseCreate(modelName, crBody.Data)

	request = CreateRequest{
		Operation: op,
	}

	return request
}

func (s CreateServer) Serve(w http.ResponseWriter, req interface{}) {
	// params := req.(CreateRequest)

	// id, ok := params.Data["id"]
	// if !ok {
	// 	panic("No id")
	// }

	// st := db.MakeStruct(params.Type)
	// err := mapstructure.Decode(params.Data, &st)
	// if err != nil {
	// 	panic(err)
	// }

	// db.DB.Insert(id, &st)

	// response := CreateResponse{Data: st}
	// bytes, _ := jsoniter.Marshal(&response)
	// _, _ = w.Write(bytes)
}
