package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
	"golang.org/x/crypto/scrypt"
)

var (
	ErrValue     = fmt.Errorf("invalid value for type")
	ErrNotStored = fmt.Errorf("value not stored")
)

func PasswordFromJSON(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrValue
	}

	value := args[0]
	password, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}

	recVal := args[1]
	rec, ok := recVal.(db.Record)
	salt, err := rec.ID().Bytes()
	if err != nil {
		return nil, err
	}

	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return dk, nil
}

var passwordValidator = db.MakeNativeFunction(
	db.MakeID("9c4a5530-c6ac-4564-b98f-9731052b881a"),
	"password",
	2,
	PasswordFromJSON,
)

var Password = db.MakeCoreDatatype(
	db.MakeID("9455e28f-2613-4d3b-bb94-34a224d49ee8"),
	"password",
	db.BytesStorage,
	passwordValidator,
)
