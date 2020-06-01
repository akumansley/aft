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
// FromJson is a reference to a Code struct
// ToJson is a reference to a Code Struct
// StorageType is the type of the raw data stored for the given datatype
// JsonType is the type of the Json for the given datatype
type Datatype struct {
	Id          uuid.UUID
	Name        string
	FromJson    Code
	ToJson      Code
	StorageType interface{}
	JsonType    interface{}
}

// bool datatype
func boolFromJsonFunc(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

// int datatype
func intFromJsonFunc(value interface{}) (interface{}, error) {
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
func enumFromJsonFunc(value interface{}) (interface{}, error) {
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
func stringFromJsonFunc(value interface{}) (interface{}, error) {
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
func textFromJsonFunc(value interface{}) (interface{}, error) {
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

func emailAddressFromJsonFunc(value interface{}) (interface{}, error) {
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
func uuidFromJsonFunc(value interface{}) (interface{}, error) {
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
func floatFromJsonFunc(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return f, nil
}

// URL datatype
// This is an example of a more complex datatype.
func urlFromJsonFunc(value interface{}) (interface{}, error) {
	urlString, ok := value.(string)
	if ok {
		u, err := url.Parse(urlString)
		if err != nil {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		} else if u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		} else if u.Scheme != "http" && u.Scheme != "https" {
			return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
		}
	} else {
		return nil, fmt.Errorf("%w: expected url got %T", ErrValue, value)
	}
	return urlString, nil
}

//Native code datatype
//This shouldn't ever get called, but is a placeholder for the db
func nativeCodeFromJsonFunc(value interface{}) (interface{}, error) {
	return nil, fmt.Errorf("%w: nativeCode does not execute type", ErrValue)
}
