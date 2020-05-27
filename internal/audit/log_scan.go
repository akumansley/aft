package audit

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
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

type LogScanHandler struct {
	log oplog.OpLog
	bus *bus.EventBus
}

func (s LogScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	var lsr LogScanRequest
	buf, _ := ioutil.ReadAll(r.Body)
	_ = jsoniter.Unmarshal(buf, &lsr)
	entries, err := s.log.Scan(lsr.Count, lsr.Offset)
	if err != nil {
		return
	}
	s.bus.Publish(lib.ParseRequest{Request: lsr})
	response := LogScanResponse{Data: entries}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
