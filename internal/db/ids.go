package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (i ID) String() string {
	u := uuid.UUID(i)
	return u.String()
}

func (l ID) Bytes() ([]byte, error) {
	u := uuid.UUID(l)
	return u.MarshalBinary()
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
