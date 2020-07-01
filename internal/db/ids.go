package db

import (
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

func MakeIDFromBytes(bytes []byte) ID {
	u, _ := uuid.FromBytes(bytes)
	return ID(u)
}

func MakeID(literal string) ID {
	return ID(uuid.MustParse(literal))
}
