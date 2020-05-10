package operations

import (
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/middleware"
	"context"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type FindOneRequestBody struct {
	Where   map[string]interface{} `json:"where"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type FindOneRequest struct {
	// TODO add Include/Select
	Operation db.FindOneOperation
	Include   db.Include
}

type FindOneResponse struct {
	Data interface{} `json:"data"`
}

type FindOneServer struct {
}

func (s FindOneServer) Parse(ctx context.Context, req *http.Request) (interface{}, error) {
	tx := middleware.TxFromContext(ctx)
	p := Parser{tx: tx}
	var foBody FindOneRequestBody
	vars := mux.Vars(req)
	modelName := vars["object"]
	buf, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(buf, &foBody)
	op, err := p.ParseFindOne(modelName, foBody.Where)
	if err != nil {
		return nil, err
	}
	inc, err := p.ParseInclude(modelName, foBody.Include)
	if err != nil {
		return nil, err
	}

	request := FindOneRequest{
		Operation: op,
		Include:   inc,
	}

	return request, nil
}

func (s FindOneServer) Serve(ctx context.Context, req interface{}) (interface{}, error) {
	tx := middleware.TxFromContext(ctx)
	params := req.(FindOneRequest)
	st, err := params.Operation.Apply(tx)
	if err != nil {
		return nil, err
	}
	responseData := params.Include.Resolve(tx, st)
	response := FindOneResponse{Data: responseData}
	return response, nil
}
