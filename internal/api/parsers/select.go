package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) consumeSelect(m db.Interface, keys api.Set, data map[string]interface{}) (operations.Select, error) {
	var s map[string]interface{}
	if v, ok := data["select"]; ok {
		s = v.(map[string]interface{})
		delete(keys, "select")
	}
	return p.parseSelect(m, s)
}

func (p Parser) parseSelection(r db.Relationship, value interface{}) (operations.Selection, error) {
	if v, ok := value.(bool); ok {
		if v {
			return operations.Selection{Relationship: r}, nil
		} else {
			return operations.Selection{}, fmt.Errorf("%w: Select specified as false", ErrInvalidStructure)
		}
	} else if args, ok := value.(map[string]interface{}); ok {
		op, err := p.parseNestedFindMany(r.Target().Name(), args)
		if err != nil {
			return operations.Selection{}, err
		}
		return operations.Selection{Relationship: r, NestedFindMany: op}, nil
	}
	return operations.Selection{}, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, value)

}

func (p Parser) parseSelect(m db.Interface, data map[string]interface{}) (s operations.Select, err error) {
	var selects []operations.Selection
	rels, err := m.Relationships()
	if err != nil {
		return
	}
	for _, rel := range rels {
		if val, ok := data[rel.Name()]; ok {
			sel, err := p.parseSelection(rel, val)
			if err != nil {
				return s, err
			}
			selects = append(selects, sel)
			delete(data, rel.Name())
		}
	}

	// delete all attributes from data
	fields := make(api.Set)
	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if value, ok := data[attr.Name()]; ok {
			if v, ok := value.(bool); ok {
				if v {
					fields[attr.Name()] = api.Void{}
					delete(data, attr.Name())
				} else {
					return s, fmt.Errorf("%w: Field specified as false", ErrInvalidStructure)
				}
			}
		}
	}
	if len(data) != 0 {
		return s, fmt.Errorf("%w: %v", ErrUnusedKeys, data)
	}
	s = operations.Select{Selecting: true, Selects: selects, Fields: fields}
	return
}
