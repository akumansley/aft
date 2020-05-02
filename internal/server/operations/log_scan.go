package operations

import (
	"awans.org/aft/internal/oplog"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type LogScanRequest struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type LogScanResponse struct {
	Data interface{} `json:"data"`
}

type LogScanServer struct {
	Log oplog.OpLog
}

func (s LogScanServer) Parse(req *http.Request) (interface{}, error) {
	var lsrBody LogScanRequest
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &lsrBody)
	return lsrBody, nil
}

func (s LogScanServer) Serve(req interface{}) (interface{}, error) {
	params := req.(LogScanRequest)
	entries, err := s.Log.Scan(params.Count, params.Offset)
	if err != nil {
		return nil, err
	}
	response := LogScanResponse{Data: entries}
	return response, nil
}
