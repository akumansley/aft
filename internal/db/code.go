package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/starlight-go/starlight"
	"reflect"
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
	AndrewFromJSON:     nil,
	AndrewToJSON:       nil,
}

func CallFunc(c Code, args interface{}, st StorageType) (interface{}, error) {
	if c.Runtime == Golang {
		return functionMap[c](args)
	} else if c.Runtime == Starlark {
		return skylarkParser(c.Code, args, st)
	}
	return nil, nil
}

type s struct {
	Value string
	Error string
}

type f struct {
	Value float64
	Error string
}

type b struct {
	Value bool
	Error string
}

//Starlark
//uses https://github.com/starlight-go/starlight
func skylarkParser(code string, args interface{}, st StorageType) (interface{}, error) {
	globals := map[string]interface{}{
		"printf": fmt.Printf,
		"errorf": fmt.Printf,
	}
	switch st {
	case BoolType:
		globals["args"] = &b{Value: args.(bool)}
	case FloatType:
		globals["args"] = &f{Value: args.(float64)}
	case StringType:
		globals["args"] = &s{Value: args.(string)}
	default:
		panic("Unrecognized storage type")
	}

	_, err := starlight.Eval([]byte(code), globals, nil)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(globals["args"]).Elem().FieldByName("Value")
	e := reflect.ValueOf(globals["args"]).Elem().FieldByName("Error")
	if e.String() != "" {
		return nil, fmt.Errorf("%s", e.String())
	}
	switch st {
	case BoolType:
		return v.Bool(), nil
	case FloatType:
		return v.Float(), nil
	case StringType:
		return v.String(), nil
	default:
		panic("Unrecognized storage type")
	}
	return nil, fmt.Errorf("Shouldn't get here in starlark")
}
