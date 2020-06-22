package db

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

var (
	ErrInvalidRelationship = fmt.Errorf("%w: invalid relationship", ErrData)
)

type ID uuid.UUID

func (l ID) Bytes() ([]byte, error) {
	u := uuid.UUID(l)
	return u.MarshalBinary()
}

func MakeID(literal string) ID {
	return ID(uuid.MustParse(literal))
}

func MakeModelID(literal string) ModelID {
	return ModelID(uuid.MustParse(literal))
}

type ModelID uuid.UUID

func (m ModelID) Bytes() ([]byte, error) {
	u := uuid.UUID(m)
	return u.MarshalBinary()
}

type Attribute struct {
	Name     string
	Datatype Datatype
	ID       ID
}

type RelType int64

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Relationship struct {
	ID           ID
	LeftBinding  RelType
	RightBinding RelType
	LeftModelID  ModelID
	RightModelID ModelID
	LeftName     string
	RightName    string
}

func (r Relationship) Left() Binding {
	return Binding{Relationship: r, Left: true}
}

func (r Relationship) Right() Binding {
	return Binding{Relationship: r, Left: false}
}

func JSONKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vID", strings.Title(strings.ToLower(key)))
}

func JSONKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

type Model struct {
	ID                 ModelID
	Name               string
	Attributes         []Attribute
	LeftRelationships  []Relationship
	RightRelationships []Relationship
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

func (m Model) GetBinding(name string) (Binding, error) {
	for _, b := range m.Bindings() {
		if b.Name() == name {
			return b, nil
		}
	}
	return Binding{}, ErrInvalidRelationship
}

func (m Model) Bindings() []Binding {
	var bs []Binding
	for _, r := range m.LeftRelationships {
		bs = append(bs, Binding{Relationship: r, Left: true})
	}
	for _, r := range m.RightRelationships {
		bs = append(bs, Binding{Relationship: r, Left: false})
	}
	return bs
}

type Binding struct {
	Relationship Relationship
	Left         bool
}

func (b Binding) HasField() bool {
	if b.Left {
		return (b.Relationship.LeftBinding == BelongsTo)
	} else {
		return (b.Relationship.RightBinding == BelongsTo)
	}
}

func (b Binding) Name() string {
	if b.Left {
		return b.Relationship.LeftName
	} else {
		return b.Relationship.RightName
	}
}

func (b Binding) ModelID() ModelID {
	if b.Left {
		return b.Relationship.LeftModelID
	} else {
		return b.Relationship.RightModelID
	}
}

func (b Binding) Dual() Binding {
	return Binding{Relationship: b.Relationship, Left: !b.Left}
}

func (b Binding) RelType() RelType {
	if b.Left {
		return b.Relationship.LeftBinding
	} else {
		return b.Relationship.RightBinding
	}
}
