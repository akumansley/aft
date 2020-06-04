package db

import (
	"github.com/google/uuid"
)

type Datatype struct {
	ID          uuid.UUID
	Name        string
	Validator   Code
	StorageType StorageType
}

type StorageType int64

const (
	BoolType StorageType = iota
	IntType
	StringType
	FloatType
	UUIDType
)

var storageType map[StorageType]interface{} = map[StorageType]interface{}{
	BoolType:   false,
	IntType:    int64(0),
	StringType: "",
	FloatType:  0.0,
	UUIDType:   uuid.UUID{},
}

var datatypeMap map[uuid.UUID]Datatype = map[uuid.UUID]Datatype{
	Bool.ID:         Bool,
	Int.ID:          Int,
	Enum.ID:         Enum,
	String.ID:       String,
	Text.ID:         Text,
	EmailAddress.ID: EmailAddress,
	UUID.ID:         UUID,
	Float.ID:        Float,
	URL.ID:          URL,
}

func (d Datatype) FromJSON(arg interface{}) (interface{}, error) {
	c := d.Validator
	return CallValidator(c, arg)
}

var Bool = Datatype{
	ID:          uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:        "bool",
	Validator:   boolValidator,
	StorageType: BoolType,
}

var Int = Datatype{
	ID:          uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:        "int",
	Validator:   intValidator,
	StorageType: IntType,
}

var Enum = Datatype{
	ID:          uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:        "enum",
	Validator:   enumValidator,
	StorageType: IntType,
}

var String = Datatype{
	ID:          uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:        "string",
	Validator:   stringValidator,
	StorageType: StringType,
}

var Text = Datatype{
	ID:          uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:        "text",
	Validator:   textValidator,
	StorageType: StringType,
}

var EmailAddress = Datatype{
	ID:          uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:        "emailAddress",
	Validator:   emailAddressValidator,
	StorageType: StringType,
}

var UUID = Datatype{
	ID:          uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:        "uuid",
	Validator:   uuidValidator,
	StorageType: UUIDType,
}

var Float = Datatype{
	ID:          uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:        "float",
	Validator:   floatValidator,
	StorageType: FloatType,
}

var URL = Datatype{
	ID:          uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:        "url",
	Validator:   URLValidator,
	StorageType: StringType,
}
