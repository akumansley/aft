package db

import (
	"fmt"
)

type mlBox struct {
	ModelL
}

type ModelL struct {
	ID         ID     `record:"id"`
	Name       string `record:"name"`
	Attributes []Attribute
}

func (lit ModelL) AsModel() Model {
	return mlBox{lit}
}

func (m mlBox) ID() ID {
	return m.ModelL.ID
}

func (m mlBox) Lit() ModelL {
	return m.ModelL
}

func (m mlBox) Name() string {
	return m.ModelL.Name
}

func (m mlBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) RelationshipByName(name string) (Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) Attributes() ([]Attribute, error) {
	return m.ModelL.Attributes, nil
}

func (m mlBox) AttributeByName(name string) (a Attribute, err error) {
	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, nil
		}
	}
	a, ok := SystemAttrs[name]
	if !ok {
		err = fmt.Errorf("No attribute on model: %v %v", m.Name, name)
	}
	return
}
