package parsers

import (
	"errors"
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
)

func (p Parser) consumeOrder(m db.Interface, keys api.Set, data map[string]interface{}) ([]db.Sort, error) {
	var c []interface{}
	if v, ok := data["order"]; ok {
		c, ok = v.([]interface{})
		if !ok {
			return nil, errors.New("Invalid structure in order; expected list")
		}
		delete(keys, "order")
	}
	return p.parseOrder(m, c)
}

// an order is
// [{"age":"desc"}, {"email":"asc"}]
func (p Parser) parseOrder(m db.Interface, list []interface{}) (order []db.Sort, err error) {
	for _, item := range list {
		object, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("Invalid structure in order")
		}
		sort, err := p.parseSort(m, object)
		if err != nil {
			return nil, err
		}
		order = append(order, sort)
	}
	return
}

// a sort is
// {"age":"desc"}
func (p Parser) parseSort(m db.Interface, obj map[string]interface{}) (sort db.Sort, err error) {
	if err != nil {
		return
	}
	if len(obj) != 1 {
		return sort, errors.New("Invalid structure in sort")
	}
	for key, val := range obj {
		_, err := m.AttributeByName(p.Tx, key)
		if err != nil {
			return sort, fmt.Errorf("sort on invalid attribute %v", key)
		}
		sort.AttributeName = key

		if val == "asc" {
			sort.Ascending = true
		} else if val == "desc" {
			sort.Ascending = false
		} else {
			return sort, fmt.Errorf("invalid direction %v", val)
		}
	}
	return
}

func (p Parser) consumeSkip(m db.Interface, keys api.Set, data map[string]interface{}) (int, error) {
	if v, ok := data["skip"]; ok {
		skipFloat := v.(float32)
		delete(keys, "skip")
		return int(skipFloat), nil
	}
	return 0, nil
}

func (p Parser) consumeTake(m db.Interface, keys api.Set, data map[string]interface{}) (int, error) {
	if v, ok := data["take"]; ok {
		takeFloat := v.(float32)
		delete(keys, "take")
		return int(takeFloat), nil
	}
	return 0, nil
}
