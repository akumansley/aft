package db

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/olebedev/go-duktape.v3"
)

type Code struct {
	Id      uuid.UUID
	Name    string
	Runtime Runtime
	Syntax  string
}

type Runtime int64

const (
	Golang Runtime = iota
	Javascript
	Starlark
)

var functionMap map[uuid.UUID]func(interface{}) (interface{}, error) = map[uuid.UUID]func(interface{}) (interface{}, error){
	boolFromJson.Id:         boolFromJsonFunc,
	intFromJson.Id:          intFromJsonFunc,
	enumFromJson.Id:         enumFromJsonFunc,
	stringFromJson.Id:       stringFromJsonFunc,
	textFromJson.Id:         textFromJsonFunc,
	emailAddressFromJson.Id: emailAddressFromJsonFunc,
	uuidFromJson.Id:         uuidFromJsonFunc,
	floatFromJson.Id:        floatFromJsonFunc,
	urlFromJson.Id:          urlFromJsonFunc,
	nativeCodeFromJson.Id:   nativeCodeFromJsonFunc,
}

func CallFunc(c Code) func(interface{}) (interface{}, error) {
	if c.Runtime == Golang {
		return functionMap[c.Id]
	}
	return nil
}

//
//Javascript datatype
//uses https://github.com/olebedev/go-duktape bindings
func javascriptFromJson(value interface{}) (interface{}, error) {
	javascriptString, ok := value.(string)
	if ok {
		ctx := duktape.New()
		err := ctx.PevalString(javascriptString)
		result := ctx.GetNumber(-1)
		ctx.DestroyHeap()
		if &err != nil {
			return result, nil
		} else {
			return nil, fmt.Errorf("%w: expected Javascript got %s", ErrValue, err)
		}
	}
	return nil, fmt.Errorf("%w: expected Javascript got %s", ErrValue, value)
}
