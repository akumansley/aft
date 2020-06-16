package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type FindManyRequestBody struct {
	Where   map[string]interface{} `json:"where"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
	OrderBy map[string]interface{} `json:"orderBy"`
	Skip    map[string]interface{} `json:"skip"`
	After   map[string]interface{} `json:"after"`
	Before  map[string]interface{} `json:"before"`
	First   map[string]interface{} `json:"first"`
	Last    map[string]interface{} `json:"last"`
}

type FindManyRequest struct {
	// TODO add Select
	Operation FindManyOperation
	Include   Include
}

type FindManyResponse struct {
	Data interface{} `json:"data"`
}

type FindManyHandler struct {
	db  db.DB
	bus *bus.EventBus
}

func (s FindManyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := s.db.NewTx()
	p := Parser{tx: tx}
	var foBody FindManyRequestBody
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	buf, _ := ioutil.ReadAll(r.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)

	// parse the request
	op, err := p.ParseFindMany(modelName, foBody.Where)
	if err != nil {
		return
	}
	inc, err := p.ParseInclude(modelName, foBody.Include)
	if err != nil {
		return
	}

	request := FindManyRequest{
		Operation: op,
		Include:   inc,
	}

	s.bus.Publish(lib.ParseRequest{Request: request})

	recs := request.Operation.Apply(tx)
	var rData []QueryResult
	for _, rec := range recs {
		responseData := request.Include.Resolve(tx, rec)
		rData = append(rData, responseData)
	}
	response := FindManyResponse{Data: rData}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
