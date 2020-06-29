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

func MakeModelID(literal string) ModelID {
	return ModelID(uuid.MustParse(literal))
}

func MakeModelIDFromBytes(bytes []byte) ModelID {
	u, _ := uuid.FromBytes(bytes)
	return ModelID(u)
}

type ModelID uuid.UUID

func (m ModelID) String() string {
	u := uuid.UUID(m)
	return u.String()
}

func (m ModelID) Bytes() ([]byte, error) {
	u := uuid.UUID(m)
	return u.MarshalBinary()
}
