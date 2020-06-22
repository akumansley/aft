package db

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
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

type Attribute struct {
	Name     string
	Datatype Datatype
	ID       ID
}

type Relationship struct {
	ID     ID
	Name   string
	Multi  bool
	Source Model
	Target Model
}

type Model struct {
	ID         ModelID
	Name       string
	Attributes []Attribute
}

func JSONKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vID", strings.Title(strings.ToLower(key)))
}

func JSONKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

func (m Model) AttributeByName(name string) Attribute {
	for _, a := range m.Attributes {
		if a.Name == name {
			return a
		}
	}
	a, ok := SystemAttrs[name]
	if !ok {
		s := fmt.Sprintf("No attribute on model: %v %v", m.Name, name)
		panic(s)
	}
	return a
}
