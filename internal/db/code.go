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

var functionMap map[Code]func(interface{}) (interface{}, error) = map[Code]func(interface{}) (interface{}, error){
	boolFromJSON:         boolFromJSONFunc,
	intFromJSON:          intFromJSONFunc,
	enumFromJSON:         enumFromJSONFunc,
	stringFromJSON:       stringFromJSONFunc,
	textFromJSON:         textFromJSONFunc,
	emailAddressFromJSON: emailAddressFromJSONFunc,
	uuidFromJSON:         uuidFromJSONFunc,
	floatFromJSON:        floatFromJSONFunc,
	URLFromJSON:          URLFromJSONFunc,
	//TODO add real ToJSON functions
	boolToJSON:         boolFromJSONFunc,
	intToJSON:          intFromJSONFunc,
	enumToJSON:         enumFromJSONFunc,
	stringToJSON:       stringFromJSONFunc,
	textToJSON:         textFromJSONFunc,
	emailAddressToJSON: emailAddressFromJSONFunc,
	uuidToJSON:         uuidFromJSONFunc,
	floatToJSON:        floatFromJSONFunc,
	URLToJSON:          URLFromJSONFunc,
}

func CallFunc(c Code, value interface{}) (interface{}, error) {
	if c.Runtime == Golang {
		return functionMap[c](value)
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
