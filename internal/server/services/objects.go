package services

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/server/db"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type InfoObjectsRequest struct {
	Id string `json:"id"`
}

type InfoObjectsResponse struct {
	Object data.Object `json:"object"`
}

type InfoObjectsServer struct{}

func (m InfoObjectsServer) Parse(req *http.Request) interface{} {
	var request InfoObjectsRequest
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &request)
	return request
}

func (m InfoObjectsServer) Serve(w http.ResponseWriter, req interface{}) {
	params := req.(InfoObjectsRequest)
	id := params.Id

	results := db.ObjectTable.Query(id)

	// if ok {
	// 	response := InfoObjectsResponse{
	// 		Object: object,
	// 	}
	bytes, _ := jsoniter.Marshal(&results)
	_, _ = w.Write(bytes)
	// }
}

type ListObjectsResponse struct {
	Objects []data.Object `json:"objects"`
}

type ListObjectsServer struct{}

func (m ListObjectsServer) Parse(req *http.Request) interface{} {
	return nil
}

func (m ListObjectsServer) Serve(w http.ResponseWriter, req interface{}) {
	resp := ListObjectsResponse{Objects: db.Objects}
	bytes, _ := jsoniter.Marshal(&resp)
	_, _ = w.Write(bytes)
}
