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

func InfoObjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var request InfoObjectsRequest
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &request)
	id := request.Id

	object, ok := db.ObjectTable.Get(id).(data.Object)

	if ok {
		response := InfoObjectsResponse{
			Object: object,
		}
		bytes, _ := jsoniter.Marshal(&response)
		_, _ = w.Write(bytes)
	} else {
		http.NotFound(w, req)
	}
}

type ListObjectsResponse struct {
	Objects []data.Object `json:"objects"`
}

func ListObjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := ListObjectsResponse{Objects: db.Objects}
	bytes, err := jsoniter.Marshal(&resp)
	if err != nil {
		http.NotFound(w, req)
	}
	_, err = w.Write(bytes)
}
