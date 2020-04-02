package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type QueryRequest struct {
	Q    string `json:"query"`
	Type string `json:-`
}

type QueryResponse struct {
	Data []interface{} `json:"data"`
}

type QueryServer struct{}

func (s QueryServer) Parse(req *http.Request) interface{} {
	var request QueryRequest
	vars := mux.Vars(req)
	type_ := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &request)
	request.Type = type_
	return request
}

func (s QueryServer) Serve(w http.ResponseWriter, req interface{}) {
	params := req.(QueryRequest)
	id := params.Q

	results := db.DB.GetTable(params.Type).Query(id)
	resp := QueryResponse{Data: results}

	bytes, _ := jsoniter.Marshal(&resp)
	_, _ = w.Write(bytes)
}
