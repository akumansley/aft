package operations

import (
	"awans.org/aft/internal/server/middleware"
	"context"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type CreateRequestBody struct {
	Data    map[string]interface{} `json:"data"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type CreateRequest struct {
	// TODO add Select
	Operation CreateOperation
	Include Include
}

type CreateResponse struct {
	Data interface{} `json:"data"`
}

type CreateServer struct {
}

func (s CreateServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	tx := middleware.RWTxFromContext(ctx)
	p := Parser{tx: tx}
	var crBody CreateRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	body, _ := ioutil.ReadAll(req.Body)
	err := jsoniter.Unmarshal(body, &crBody)
	if err != nil {
		return nil, err
	}
	
	var request CreateRequest
	op, err := p.ParseCreate(modelName, crBody.Data)
	if err != nil {
		return nil, err
	}
	inc, err := p.ParseInclude(modelName, crBody.Include)
	if err != nil {
		return nil, err
	}
	
	request = CreateRequest{
		Operation: op,
		Include: inc,
	}
	return request, nil
}

func (s CreateServer) Serve(ctx context.Context, req interface{}) (interface{}, error) {
	tx := middleware.RWTxFromContext(ctx)
	params := req.(CreateRequest)
	st, err := params.Operation.Apply(tx)
	if err != nil {
		return nil, err
	}
	responseData := params.Include.Resolve(tx, st)
	response := CreateResponse{Data: responseData}
	return response, nil
}
