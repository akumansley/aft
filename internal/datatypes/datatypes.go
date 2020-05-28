package datatypes

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"regexp"
)

var (
	ErrData        = errors.New("data-error")
	ErrInvalidAttr = fmt.Errorf("%w: invalid attribute", ErrData)
	ErrValue       = fmt.Errorf("%w: invalid value for type", ErrData)
)

type AttrType int64

const (
	Int AttrType = iota
	String
	Text
	Float
	Enum
	UUID
	Bool
	EmailAddress
)

type Datatype interface {
	fromJson(interface{}) (Datatype, error)
	ToJson() interface{}
}

func Parse(t AttrType, value interface{}) (Datatype, error) {
	switch t {
	case Int:
		i := integer{}
		return i.fromJson(value)
	case String:
		s := stringer{}
		return s.fromJson(value)
	case Text:
		t := text{}
		return t.fromJson(value)
	case Float:
		f := floating{}
		return f.fromJson(value)
	case Enum:
		e := enum{}
		return e.fromJson(value)
	case UUID:
		u := internaluuid{}
		return u.fromJson(value)
	case Bool:
		b := boolean{}
		return b.fromJson(value)
	case EmailAddress:
		e := emailAddress{}
		return e.fromJson(value)
	}
	return nil, fmt.Errorf("%w: got attribute type %v", ErrInvalidAttr, t)
}

//booleans
type boolean struct {
	value bool
}

func (this boolean) fromJson(value interface{}) (Datatype, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return boolean{value: b}, nil
}

func (this boolean) ToJson() (v interface{}) {
	return this.value
}

//Integers
type integer struct {
	value int64
}

func (this integer) fromJson(value interface{}) (Datatype, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return integer{value: i}, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return integer{value: i}, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return integer{value: i64Val}, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

func (this integer) ToJson() (v interface{}) {
	return this.value
}

//Enum
type enum struct {
	value int64
}

func (this enum) fromJson(value interface{}) (Datatype, error) {
	f, ok := value.(float64)
	if ok {
		i := int64(f)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return integer{value: i}, nil
	}
	intVal, ok := value.(int)
	if ok {
		i := int64(intVal)
		if !ok {
			return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
		}
		return integer{value: i}, nil
	}
	i64Val, ok := value.(int64)
	if ok {
		return enum{value: i64Val}, nil
	} else {
		return nil, fmt.Errorf("%w: expected int/enum got %T", ErrValue, value)
	}
}

func (this enum) ToJson() (v interface{}) {
	return this.value
}

//String
type stringer struct {
	value string
}

func (this stringer) fromJson(value interface{}) (Datatype, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string/text got %T", ErrValue, value)
	}
	return stringer{value: s}, nil
}

func (this stringer) ToJson() (v interface{}) {
	return this.value
}

//Text
type text struct {
	value string
}

func (this text) fromJson(value interface{}) (Datatype, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string/text got %T", ErrValue, value)
	}
	return text{value: s}, nil
}

func (this text) ToJson() (v interface{}) {
	return this.value
}

//EmailAddress
type emailAddress struct {
	value string
}

//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func (this emailAddress) fromJson(value interface{}) (Datatype, error) {
	emailAddressString, ok := value.(string)
	if ok {
		if (len(emailAddressString) > 254 || !matchEmail(emailAddressString)) && len(emailAddressString) != 0 {
			return nil, fmt.Errorf("expected email address got %v", emailAddressString)
		}
	} else {
		return nil, fmt.Errorf("%w: expected email address got %T", ErrValue, value)
	}
	return emailAddress{value: emailAddressString}, nil
}

func (this emailAddress) ToJson() (v interface{}) {
	return this.value
}

//UUID
type internaluuid struct {
	value uuid.UUID
}

func (this internaluuid) fromJson(value interface{}) (Datatype, error) {
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
	return internaluuid{value: u}, nil
}

func (this internaluuid) ToJson() (v interface{}) {
	return this.value
}

//Float
type floating struct {
	value float64
}

func (this floating) fromJson(value interface{}) (Datatype, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return floating{value: f}, nil
}

func (this floating) ToJson() (v interface{}) {
	return this.value
}
