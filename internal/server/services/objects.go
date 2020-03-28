package services

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/server/db"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"net/http"
)

func InfoObject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var params InfoObjectsRequest
	buf, _ := ioutil.ReadAll(req.Body)
	_ = protojson.Unmarshal(buf, &params)
	id := params.Id

	object := db.ObjectTable.Get(id).(*data.Object)
	if object != nil {
		response := InfoObjectsResponse{
			Object: object,
		}
		bytes, _ := protojson.Marshal(&response)
		_, _ = w.Write(bytes)
	} else {
		http.NotFound(w, req)
	}
}

func ListObjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := ListObjectsResponse{Data: db.Objects}
	bytes, err := protojson.Marshal(&resp)
	if err != nil {
		http.NotFound(w, req)
	}
	_, err = w.Write(bytes)
}
