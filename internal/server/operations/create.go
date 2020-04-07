package operations

import (
	"awans.org/aft/internal/server/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"reflect"
)

type CreateRequestBody struct {
	Data    map[string]interface{} `json:"data"`
	Select  map[string]interface{} `json:"select"`
	Include map[string]interface{} `json:"include"`
}

type CreateRequest struct {
	// TODO add Include/Select
	Operation CreateOperation
}

type CreateResponse struct {
	Data interface{} `json:"data"`
}

type CreateOperation struct {
	Struct interface{}
	Nested []NestedCreateOperation
}

// todo do we need separate structs for each nested variant or can we
// somehow reuse the raw ops and just parse into them with the added metadata
type NestedCreateOperation struct {
	Parent       interface{}
	Relationship string
	Struct       interface{}
	Nested       []NestedCreateOperation
}

type NestedConnectOperation struct {
	Parent       interface{}
	Relationship string
	Struct       interface{}
	Nested       []NestedCreateOperation
}

type CreateServer struct{}

func (s CreateServer) Parse(req *http.Request) interface{} {
	var crBody CreateRequestBody
	vars := mux.Vars(req)
	type_ := vars["object"]
	body, _ := ioutil.ReadAll(req.Body)
	_ = jsoniter.Unmarshal(body, &crBody)
	var request CreateRequest
	var rootOp CreateOperation
	st := db.MakeStruct(type_)

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &st,
		WeaklyTypedInput: true,
		DecodeHook: func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
			if from == reflect.TypeOf("") && to == reflect.TypeOf(uuid.UUID{}) {
				idString, ok := data.(string)
				if !ok {
					return nil, errors.New("data was not a string")
				}

				u, err := uuid.Parse(idString)
				if err != nil {
					return nil, err
				}
				return u, nil
			}
			return data, nil
		},
	})
	fmt.Printf("Calling decode!!")

	decoder.Decode(crBody.Data)

	rootOp = CreateOperation{
		Struct: st,
		Nested: []NestedCreateOperation{},
	}

	request = CreateRequest{
		Operation: rootOp,
	}

	return request
}

func (s CreateServer) Serve(w http.ResponseWriter, req interface{}) {
	// params := req.(CreateRequest)

	// id, ok := params.Data["id"]
	// if !ok {
	// 	panic("No id")
	// }

	// st := db.MakeStruct(params.Type)
	// err := mapstructure.Decode(params.Data, &st)
	// if err != nil {
	// 	panic(err)
	// }

	// db.DB.Insert(id, &st)

	// response := CreateResponse{Data: st}
	// bytes, _ := jsoniter.Marshal(&response)
	// _, _ = w.Write(bytes)
}
