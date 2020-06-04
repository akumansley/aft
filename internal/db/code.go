package db

import (
	"awans.org/aft/internal/runtime"
	"github.com/google/uuid"
	"fmt"
)

type Code struct {
	ID       uuid.UUID
	Name     string
	Function Function
	Runtime  Runtime
	Code     string
}

type Function int64

const (
	Validator Function = iota
)

type Runtime int64

const (
	Golang Runtime = iota
	Javascript
	Starlark
)

var codeMap map[uuid.UUID]Code = map[uuid.UUID]Code{
	boolValidator.ID:         boolValidator,
	intValidator.ID:          intValidator,
	enumValidator.ID:         enumValidator,
	stringValidator.ID:       stringValidator,
	textValidator.ID:         textValidator,
	emailAddressValidator.ID: emailAddressValidator,
	uuidValidator.ID:         uuidValidator,
	floatValidator.ID:        floatValidator,
	URLValidator.ID:          URLValidator,
}

func CallValidator(c Code, arg interface{}) (interface{}, error) {
	var fh runtime.FunctionHandle
	if(c.Runtime == Golang) {
		fh = &runtime.GoFunctionHandle{Name: c.Name}
	} else if(c.Runtime == Starlark) {
		fh = &runtime.StarlarkFunctionHandle{Code: c.Code}
	} else {
		return nil, fmt.Errorf("Unrecognized runtime")
	}
	return fh.Invoke(arg)
}

var boolValidator = Code{
	ID:       uuid.MustParse("8e806967-c462-47af-8756-48674537a909"),
	Name:     "boolValidator",
	Runtime:  Golang,
	Function: Validator,
}

var intValidator = Code{
	ID:       uuid.MustParse("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Name:     "intValidator",
	Runtime:  Golang,
	Function: Validator,
}

var enumValidator = Code{
	ID:       uuid.MustParse("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
	Name:     "enumValidator",
	Runtime:  Golang,
	Function: Validator,
}

var stringValidator = Code{
	ID:       uuid.MustParse("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Name:     "stringValidator",
	Runtime:  Golang,
	Function: Validator,
}

var textValidator = Code{
	ID:       uuid.MustParse("9f10ac9f-afd2-423a-8857-d900a0c97563"),
	Name:     "textValidator",
	Runtime:  Golang,
	Function: Validator,
}

var emailAddressValidator = Code{
	ID:       uuid.MustParse("ed046b08-ade2-4570-ade4-dd1e31078219"),
	Name:     "emailAddressValidator",
	Runtime:  Golang,
	Function: Validator,
}

var uuidValidator = Code{
	ID:       uuid.MustParse("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Name:     "uuidValidator",
	Runtime:  Golang,
	Function: Validator,
}

var floatValidator = Code{
	ID:       uuid.MustParse("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Name:     "floatValidator",
	Runtime:  Golang,
	Function: Validator,
}

var URLValidator = Code{
	ID:       uuid.MustParse("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	Name:     "urlValidator",
	Runtime:  Golang,
	Function: Validator,
}
