package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) parseInclusion(r db.Relationship, value interface{}) operations.Inclusion {
	if v, ok := value.(bool); ok {
		if v {
			return operations.Inclusion{Relationship: r}
		} else {
			panic("Include specified as false?")
		}
	}
	panic("Include with findMany args not yet implemented")
}

func (p Parser) ParseInclude(modelName string, data map[string]interface{}) (i operations.Include, err error) {
	fmt.Println(modelName)
	fmt.Println(data)
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}
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
		inc := p.parseInclusion(r, val)
		includes = append(includes, inc)
	}
	i = operations.Include{Includes: includes}
	return
}
