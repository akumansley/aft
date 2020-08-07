package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) consumeInclude(m db.Interface, keys api.Set, data map[string]interface{}) (operations.Include, error) {
	var i map[string]interface{}
	if v, ok := data["include"]; ok {
		i = v.(map[string]interface{})
		delete(keys, "include")
	}
	return p.parseInclude(m, i)
}

func (p Parser) parseInclusion(r db.Relationship, value interface{}) (operations.Inclusion, error) {
	if v, ok := value.(bool); ok {
		if v {
			return operations.Inclusion{Relationship: r}, nil
		} else {
			return operations.Inclusion{}, fmt.Errorf("%w: Include specified as false", ErrInvalidStructure)
		}
	} else if args, ok := value.(map[string]interface{}); ok {
		op, err := p.parseNestedFindMany(r.Target().Name(), args)
		if err != nil {
			return operations.Inclusion{}, err
		}
		return operations.Inclusion{Relationship: r, NestedFindMany: op}, nil
	}
	return operations.Inclusion{}, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, value)

}

func (p Parser) parseInclude(m db.Interface, data map[string]interface{}) (i operations.Include, err error) {
	var includes []operations.Inclusion
	rels, err := m.Relationships()
	relsByName := map[string]db.Relationship{}
	for _, r := range rels {
		relsByName[r.Name()] = r
	}

	for k, val := range data {
		r, ok := relsByName[k]
		if !ok {
			err = fmt.Errorf("%w: %v", ErrInvalidRelationship, k)
			return
		}
		inc, err := p.parseInclusion(r, val)
		if err != nil {
			return i, err
		}
		includes = append(includes, inc)
	}
	i = operations.Include{Includes: includes}
	return
}
