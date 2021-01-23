package db

import (
	"encoding/gob"
	"encoding/json"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (i ID) String() string {
	u := uuid.UUID(i)
	return u.String()
}

func (l ID) Bytes() []byte {
	u := uuid.UUID(l)
	v, _ := u.MarshalBinary()
	return v
}

func (l ID) MarshalJSON() ([]byte, error) {
	u := uuid.UUID(l)
	return json.Marshal(u.String())
}

func MakeIDFromBytes(bytes []byte) ID {
	u, _ := uuid.FromBytes(bytes)
	return ID(u)
}

func MakeID(literal string) ID {
	return ID(uuid.MustParse(literal))
}

func NewID() ID {
	return ID(uuid.New())
}

func init() {
	gob.Register(ID{})
	gob.Register([]interface{}{})
}
