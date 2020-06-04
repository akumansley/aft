package runtime

import (
	"awans.org/aft/internal/datatypes"
	"fmt"
	"github.com/starlight-go/starlight"
	"reflect"
)

type FunctionHandle interface{ 
	Invoke(arg interface{}) (interface{}, error)
}

type GoFunctionHandle struct {
	Name string
}

type StarlarkFunctionHandle struct {
	Code string
}

func (g *GoFunctionHandle) Invoke(arg interface{}) (interface{}, error){
	if f, ok := datatypes.FunctionMap[g.Name]; ok {
		return f(arg)
	}
	return nil, fmt.Errorf("Func not found")
}

//Starlark
//uses https://github.com/starlight-go/starlight
func (s *StarlarkFunctionHandle) Invoke(arg interface{}) (interface{}, error){
	globals := map[string]interface{}{
		"printf": fmt.Printf,
	}
	switch arg.(type) {
	case int64:
		type i struct {
			Value int64
			Error string
		}
		globals["arg"] = &i{Value: arg.(int64)}
	case bool:
		type b struct {
			Value bool
			Error string
		}
		globals["arg"] = &b{Value: arg.(bool)}
	case float64:
		type f struct {
			Value float64
			Error string
		}
		globals["arg"] = &f{Value: arg.(float64)}
	case string:
		type s struct {
			Value string
			Error string
		}
		globals["arg"] = &s{Value: arg.(string)}
	default:
		panic("Unrecognized storage format")
	}

	_, err := starlight.Eval([]byte(s.Code), globals, nil)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(globals["arg"]).Elem().FieldByName("Value")
	e := reflect.ValueOf(globals["arg"]).Elem().FieldByName("Error")
	if e.String() != "" {
		return nil, fmt.Errorf("%s", e.String())
	}
	switch arg.(type) {
	case int64:
		return v.Int(), nil
	case bool:
		return v.Bool(), nil
	case float64:
		return v.Float(), nil
	case string:
		return v.String(), nil
	default:
		panic("Unrecognized storage type")
	}
	return nil, fmt.Errorf("Shouldn't get here in starlark")
}
