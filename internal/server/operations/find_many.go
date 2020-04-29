package operations

import (
	"awans.org/aft/internal/server/db"
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
	// TODO add Include/Select
	Operation db.FindManyOperation
	Include   db.Include
}

type FindManyResponse struct {
	Data interface{} `json:"data"`
}

type FindManyServer struct {
	DB db.DB
}

func (s FindManyServer) Parse(req *http.Request) (interface{}, error) {
	p := Parser{db: s.DB}
	var foBody FindManyRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)
	op := p.ParseFindMany(modelName, foBody.Where)
	inc := p.ParseInclude(modelName, foBody.Include)

	request := FindManyRequest{
		Operation: op,
		Include:   inc,
	}

	return request, nil
}

func (s FindManyServer) Serve(req interface{}) (interface{}, error) {
	params := req.(FindManyRequest)
	stIf := params.Operation.Apply(s.DB)
	stLs := stIf.([]interface{})
	var rData []interface{}
	for _, st := range stLs {
		responseData := params.Include.Resolve(s.DB, st)
		rData = append(rData, responseData)
	}
	response := FindManyResponse{Data: rData}
	return response, nil
}
