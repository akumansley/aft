package operations

import (
	"awans.org/aft/internal/server/middleware"
	"context"
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
	Operation FindManyOperation
	Include   Include
}

type FindManyResponse struct {
	Data interface{} `json:"data"`
}

type FindManyServer struct {
}

func (s FindManyServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	tx := middleware.TxFromContext(ctx)
	p := Parser{tx: tx}
	var foBody FindManyRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)
	op, err := p.ParseFindMany(modelName, foBody.Where)
	if err != nil {
		return nil, err
	}
	inc, err := p.ParseInclude(modelName, foBody.Include)
	if err != nil {
		return nil, err
	}

	request := FindManyRequest{
		Operation: op,
		Include:   inc,
	}

	return request, nil
}

func (s FindManyServer) Serve(ctx context.Context, req interface{}) (interface{}, error) {
	tx := middleware.TxFromContext(ctx)
	params := req.(FindManyRequest)
	recs := params.Operation.Apply(tx)
	var rData []IncludeResult
	for _, rec := range recs {
		responseData := params.Include.Resolve(tx, rec)
		rData = append(rData, responseData)
	}
	response := FindManyResponse{Data: rData}
	return response, nil
}
