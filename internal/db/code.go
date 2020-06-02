package db

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/olebedev/go-duktape.v3"
)

type Code struct {
	ID       uuid.UUID
	Name     string
	Function Function
	Runtime  Runtime
	Code     string
}

type Runtime int64

const (
	Golang Runtime = iota
	Javascript
	Starlark
)

type Function int64

const (
	FromJson Function = iota
	ToJson
)

var functionMap map[uuid.UUID]func(interface{}) (interface{}, error) = map[uuid.UUID]func(interface{}) (interface{}, error){
	boolFromJson.ID:         boolFromJsonFunc,
	intFromJson.ID:          intFromJsonFunc,
	enumFromJson.ID:         enumFromJsonFunc,
	stringFromJson.ID:       stringFromJsonFunc,
	textFromJson.ID:         textFromJsonFunc,
	emailAddressFromJson.ID: emailAddressFromJsonFunc,
	uuidFromJson.ID:         uuidFromJsonFunc,
	floatFromJson.ID:        floatFromJsonFunc,
	urlFromJson.ID:          urlFromJsonFunc,
}

func CallFunc(c Code, value interface{}) (interface{}, error) {
	if c.Runtime == Golang {
		return functionMap[c.ID](value)
	} else if c.Runtime == Javascript {
		javascriptParser(c.Code, value)
	}
	return nil, nil
}

//
//Javascript
//uses https://github.com/olebedev/go-duktape bindings
func javascriptParser(syntax string, value interface{}) (interface{}, error) {
	ctx := duktape.New()
	err := ctx.PevalString(syntax)
	result := ctx.GetNumber(-1)
	ctx.DestroyHeap()
	if &err != nil {
		return result, nil
	}
	return nil, fmt.Errorf("%w: expected Javascript got %s", ErrValue, err)
}
