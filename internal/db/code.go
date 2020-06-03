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

func CallFunc(c Code, sf StorageFormat, arg interface{}) (interface{}, error) {
	if c.Runtime == Golang {
		if f, ok := nativeFunctionMap[c]; ok {
			return f(arg)
		} else {
			return nil, fmt.Errorf("Func %s not found in native functions", c.Name)
		}
	} else if c.Runtime == Starlark {
		return skylarkParser(c.Code, sf, arg)
	}
	return nil, nil
}

//Starlark
//uses https://github.com/starlight-go/starlight
func skylarkParser(code string, sf StorageFormat, arg interface{}) (interface{}, error) {
	globals := map[string]interface{}{
		"printf": fmt.Printf,
	}
	switch sf {
	case IntFormat:
		type i struct {
			Value int64
			Error string
		}
		globals["arg"] = &i{Value: arg.(int64)}
	case BoolFormat:
		type b struct {
			Value bool
			Error string
		}
		globals["arg"] = &b{Value: arg.(bool)}
	case FloatFormat:
		type f struct {
			Value float64
			Error string
		}
		globals["arg"] = &f{Value: arg.(float64)}
	case StringFormat:
		type s struct {
			Value string
			Error string
		}
		globals["arg"] = &s{Value: arg.(string)}
	default:
		panic("Unrecognized storage format")
	}

	_, err := starlight.Eval([]byte(code), globals, nil)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(globals["arg"]).Elem().FieldByName("Value")
	e := reflect.ValueOf(globals["arg"]).Elem().FieldByName("Error")
	if e.String() != "" {
		return nil, fmt.Errorf("%s", e.String())
	}
	switch sf {
	case IntFormat:
		return v.Int(), nil
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
