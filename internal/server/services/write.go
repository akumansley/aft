package services

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type CreateRequest struct {
	Body map[string]interface{} `json:"body"`
	Type string                 `json:-`
}

type CreateResponse struct {
	Body map[string]interface{} `json:"body"`
}

type CreateServer struct{}

type Entity struct {
	Body map[string]interface{}
}

func (e Entity) GetId() string {
	val, _ := e.Body["id"].(string)
	return val
}

func (s CreateServer) Parse(req *http.Request) interface{} {
	var request CreateRequest
	vars := mux.Vars(req)
	type_ := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &request)
	request.Type = type_
	return request
}

func (s CreateServer) Serve(w http.ResponseWriter, req interface{}) {
	params := req.(CreateRequest)
	body := params.Body
	entity := Entity{Body: body}

	db.DB.GetTable(params.Type).Put(entity)
	response := CreateResponse{Body: entity.Body}
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
}
