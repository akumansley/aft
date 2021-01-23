package auth

import (
	"bytes"
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
)

var (
	ErrValue     = fmt.Errorf("invalid value for type")
	ErrNotStored = fmt.Errorf("value not stored")
)

func checkPassword(args []interface{}) (result interface{}, err error) {
	password := args[0].(string) // password
	id := args[1].(uuid.UUID)    // salt
	stored := args[2].([]byte)   // stored

	salt, _ := id.MarshalBinary()
	hashed, err := doHash([]byte(password), salt)
	return bytes.Equal(hashed, stored), err
}

var CheckPassword = db.MakeNativeFunction(
	db.MakeID("9efc2295-a143-4e76-a6e4-02a526348686"),
	"checkPassword",
	3,
	checkPassword,
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
	salt := rec.ID().Bytes()

	return doHash([]byte(password), salt)
}

func hashPassword(id db.ID, pw string) ([]byte, error) {
	return doHash([]byte(pw), id.Bytes())
}

func doHash(password, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 1<<15, 8, 1, 32)
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
