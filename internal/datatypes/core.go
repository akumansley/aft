package datatypes

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

var (
	ErrValue = fmt.Errorf("invalid value for type")
)

type GoFunctionHandle struct {
	Function func(interface{}) (interface{}, error)
}

func (g *GoFunctionHandle) Invoke(arg interface{}) (interface{}, error) {
	if g.Function != nil {
		return g.Function(arg)
	}
	return nil, fmt.Errorf("Go func not found")
}

func BoolFromJSON(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

func IntFromJSON(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int64:
		return value, nil
	case int:
		return int64(value.(int)), nil
	case float64:
		return int64(value.(float64)), nil
	case string:
		out, err := strconv.Atoi(value.(string))
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("%w: expected int got %T", ErrValue, value)

}

func StringFromJSON(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func UUIDFromJSON(value interface{}) (interface{}, error) {
	var u uuid.UUID
	var err error
	switch value.(type) {
	case string:
		u, err = uuid.Parse(value.(string))
		if err != nil {
			return nil, fmt.Errorf("%w: expected uuid got %v", ErrValue, err)
		}
	case uuid.UUID:
		u = value.(uuid.UUID)
	default:
		return nil, fmt.Errorf("%w: expected uuid got %T", ErrValue, value)
	}
	return u, nil
}

func FloatFromJSON(value interface{}) (interface{}, error) {
	switch value.(type) {
	case float64:
		return value, nil
	case int64:
		return float64(value.(int64)), nil
	case int:
		return float64(value.(int)), nil
	case string:
		out, err := strconv.ParseFloat(value.(string), 64)
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
}
