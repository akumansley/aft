package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type WriteRequest struct {
	Body map[string]interface{} `json:"body"`
	Type string                 `json:-`
}

type WriteResponse struct {
	Body map[string]interface{} `json:"body"`
}

type WriteServer struct{}

type Entity struct {
	Body map[string]interface{}
}

func (e Entity) GetId() string {
	val, _ := e.Body["id"].(string)
	return val
}

func (s WriteServer) Parse(req *http.Request) interface{} {
	var request WriteRequest
	vars := mux.Vars(req)
	type_ := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &request)
	request.Type = type_
	return request
}

func (s WriteServer) Serve(w http.ResponseWriter, req interface{}) {
	params := req.(WriteRequest)
	body := params.Body
	entity := Entity{Body: body}

	// write to the db
	db.DB.GetTable(params.Type).Put(entity)

	response := WriteResponse{Body: entity.Body}
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
}
