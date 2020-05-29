package datatypes

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"net/url"
)

var (
	ErrData        = errors.New("data-error")
	ErrValue       = fmt.Errorf("%w: invalid value for type", ErrData)
)

type Datatype interface {
	FromJson(interface{}) (interface{}, error)
	Marshal() int64
	Type() interface{}
}

var Bool = boolean{}
var Integer = integer{}
var Enum = enum{}
var String = stringer{}
var Text = text{}
var EmailAddress = emailAddress{}
var UUID = internaluuid{}
var Float = floating{}
var URL = internalurl{}

func Unmarshal(i int64) Datatype {
	switch i{
	case Bool.Marshal():
		return Bool
	case Integer.Marshal():
		return Integer
	case Enum.Marshal():
		return Enum
	case String.Marshal():
		return String
	case Text.Marshal():
		return Text
	case EmailAddress.Marshal():
		return EmailAddress
	case UUID.Marshal():
		return UUID
	case Float.Marshal():
		return Float
	case URL.Marshal():
		return URL
	}
	return nil
}

//booleans
type boolean struct{}

func (this boolean) FromJson(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

func (this boolean) Marshal() int64 {
	return 0
}

func (this boolean) Type() interface{} {
	return false
}

//Integers
type integer struct{}

func (this integer) FromJson(value interface{}) (interface{}, error) {
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

func (this integer) Marshal() int64 {
	return 1
}

func (this integer) Type() interface{} {
	return int64(0)
}

//Enum
type enum struct{}

func (this enum) FromJson(value interface{}) (interface{}, error) {
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

func (this enum) Marshal() int64 {
	return 2
}

func (this enum) Type() interface{} {
	return int64(0)
}

//String
type stringer struct{}

func (this stringer) FromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func (this stringer) Marshal() int64 {
	return 3
}

func (this stringer) Type() interface{} {
	return ""
}

//Text
type text struct{}

func (this text) FromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected text got %T", ErrValue, value)
	}
	return s, nil
}

func (this text) Marshal() int64 {
	return 4
}

func (this text) Type() interface{} {
	return ""
}

//EmailAddress
type emailAddress struct{}

//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func (this emailAddress) FromJson(value interface{}) (interface{}, error) {
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

func (this emailAddress) Marshal() int64 {
	return 5
}

func (this emailAddress) Type() interface{} {
	return ""
}

//UUID
type internaluuid struct{}

func (this internaluuid) FromJson(value interface{}) (interface{}, error) {
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

func (this internaluuid) Marshal() int64 {
	return 6
}

func (this internaluuid) Type() interface{} {
	return uuid.UUID{}
}

//Float
type floating struct{}

func (this floating) FromJson(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return f, nil
}

func (this floating) Marshal() int64 {
	return 7
}

func (this floating) Type() interface{} {
	return 0.0
}

//URL
type internalurl struct{}

func (this internalurl) FromJson(value interface{}) (interface{}, error) {
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

func (this internalurl) Marshal() int64 {
	return 8
}

func (this internalurl) Type() interface{} {
	return ""
}

