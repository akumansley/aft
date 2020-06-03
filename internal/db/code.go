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
	Validator Function = iota
)

var nativeFunctionMap map[Code]func(interface{}) (interface{}, error) = map[Code]func(interface{}) (interface{}, error){
	boolValidator:         boolValidatorFunc,
	intValidator:          intValidatorFunc,
	enumValidator:         enumValidatorFunc,
	stringValidator:       stringValidatorFunc,
	textValidator:         textValidatorFunc,
	emailAddressValidator: emailAddressValidatorFunc,
	uuidValidator:         uuidValidatorFunc,
	floatValidator:        floatValidatorFunc,
	URLValidator:          URLValidatorFunc,
}

func CallFunc(c Code, sf StorageFormat, args interface{}) (interface{}, error) {
	if c.Runtime == Golang {
		return nativeFunctionMap[c](args)
	} else if c.Runtime == Starlark {
		return skylarkParser(c.Code, sf, args)
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
func skylarkParser(code string, sf StorageFormat, args interface{}) (interface{}, error) {
	globals := map[string]interface{}{
		"printf": fmt.Printf,
		"errorf": fmt.Printf,
	}
	switch sf {
	case BoolFormat:
		globals["args"] = &b{Value: args.(bool)}
	case FloatFormat:
		globals["args"] = &f{Value: args.(float64)}
	case StringFormat:
		globals["args"] = &s{Value: args.(string)}
	default:
		panic("Unrecognized storage format")
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
	switch sf {
	case BoolFormat:
		return v.Bool(), nil
	case FloatFormat:
		return v.Float(), nil
	case StringFormat:
		return v.String(), nil
	default:
		panic("Unrecognized storage type")
	}
	return nil, fmt.Errorf("Shouldn't get here in starlark")
}
