package db

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"regexp"
)

var (
	ErrValue = fmt.Errorf("%w: invalid value for type", ErrData)
)

// Id is the UUID of the datatype
// Name is the plain english name of the type
// FromJSON is a reference to a Code struct
// ToJSON is a reference to a Code Struct
// Type is the enum of types for data storage
type Datatype struct {
	ID       uuid.UUID
	Name     string
	FromJSON Code
	ToJSON   Code
	Type     Type
}

type Type int64

const (
	BoolType Type = iota
	IntType
	StringType
	FloatType
	UUIDType
)

var storageTypeMap map[Type]interface{} = map[Type]interface{}{
	BoolType:   false,
	IntType:    int64(0),
	StringType: "",
	FloatType:  0.0,
	UUIDType:   uuid.UUID{},
}

var jsonTypeMap map[Type]interface{} = map[Type]interface{}{
	BoolType:   false,
	IntType:    0.0,
	StringType: "",
	FloatType:  0.0,
	UUIDType:   uuid.UUID{},
}

// bool datatype
func boolFromJSONFunc(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

// int datatype
func intFromJSONFunc(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return i64Val, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

// enum datatype
func enumFromJSONFunc(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return i, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return i64Val, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

// string datatype
func stringFromJSONFunc(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func stringType() interface{} {
	return ""
}

// text datatype
func textFromJSONFunc(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected text got %T", ErrValue, value)
	}
	return s, nil
}

// Email Address datatype.
// This is an example of a more complex datatype
//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func emailAddressFromJSONFunc(value interface{}) (interface{}, error) {
	emailAddressString, ok := value.(string)
	if ok {
		if (len(emailAddressString) > 254 || !matchEmail(emailAddressString)) && len(emailAddressString) != 0 {
			return nil, fmt.Errorf("%w: expected email address got %v", ErrValue, emailAddressString)
		}
	} else {
		return nil, fmt.Errorf("%w: expected email address got %T", ErrValue, value)
	}
	return emailAddressString, nil
}

// UUID datatype. Uses UUID from google underneath
func uuidFromJSONFunc(value interface{}) (interface{}, error) {
	var u uuid.UUID
	uuidString, ok := value.(string)
	if ok {
		var err error
		u, err = uuid.Parse(uuidString)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrValue, err)
		}

	} else {
		u, ok = value.(uuid.UUID)
		if !ok {
			return nil, fmt.Errorf("%w: expected uuid got %T", ErrValue, value)
		}
	}
	return u, nil
}

// float datatype.
func floatFromJSONFunc(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return f, nil
}

// URL datatype
// This is an example of a more complex datatype.
func URLFromJSONFunc(value interface{}) (interface{}, error) {
	URLString, ok := value.(string)
	if ok {
		u, err := url.Parse(URLString)
		if err != nil {
			return nil, fmt.Errorf("%w: expected URL got %T", ErrValue, value)
		} else if u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("%w: expected URL got %T", ErrValue, value)
		} else if u.Scheme != "http" && u.Scheme != "https" {
			return nil, fmt.Errorf("%w: expected URL got %T", ErrValue, value)
		}
	} else {
		return nil, fmt.Errorf("%w: expected URL got %T", ErrValue, value)
	}
	return URLString, nil
}
