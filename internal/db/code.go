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
	FromJSON Function = iota
	ToJSON
)

var functionMap map[uuid.UUID]func(interface{}) (interface{}, error) = map[uuid.UUID]func(interface{}) (interface{}, error){
	boolFromJSON.ID:         boolFromJSONFunc,
	intFromJSON.ID:          intFromJSONFunc,
	enumFromJSON.ID:         enumFromJSONFunc,
	stringFromJSON.ID:       stringFromJSONFunc,
	textFromJSON.ID:         textFromJSONFunc,
	emailAddressFromJSON.ID: emailAddressFromJSONFunc,
	uuidFromJSON.ID:         uuidFromJSONFunc,
	floatFromJSON.ID:        floatFromJSONFunc,
	URLFromJSON.ID:          URLFromJSONFunc,
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
