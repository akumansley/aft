package db

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"regexp"
)

var (
	ErrInvalidAttr         = fmt.Errorf("%w: invalid attribute", ErrData)
	ErrInvalidRelationship = fmt.Errorf("%w: invalid relationship", ErrData)
	ErrValue               = fmt.Errorf("%w: invalid value for type", ErrData)
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

type Attribute struct {
	AttrType AttrType
	Id       uuid.UUID
}

func (a Attribute) ParseFromJson(value interface{}) (interface{}, error) {
	switch a.AttrType {
	case Bool:
		b, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
		}
		return b, nil
	case Int, Enum:
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
	case String, Text:
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: expected string/text got %T", ErrValue, value)
		}
		return s, nil
	case EmailAddress:
		emailAddressString, ok := value.(string)
	    if ok {
	    	//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
			var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if (len(emailAddressString) > 254 || !rxEmail.MatchString(emailAddressString)) && len(emailAddressString) != 0 {
				return nil, fmt.Errorf("expected email address got %v", emailAddressString)
			}
		} else {
			return nil, fmt.Errorf("%w: expected email address got %T", ErrValue, value)
		}
		return emailAddressString, nil
	case Float:
		f, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
		}
		return f, nil
	case UUID:
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
	return nil, fmt.Errorf("%w: got attribute type %v", ErrInvalidAttr, a.AttrType)
}

// arguably this belongs outside of the struct
func (a Attribute) SetField(name string, value interface{}, st interface{}) error {
	fieldName := JsonKeyToFieldName(name)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	parsedValue, err := a.ParseFromJson(value)
	if err != nil {
		return err
	}
	switch parsedValue.(type) {
	case bool:
		b := parsedValue.(bool)
		field.SetBool(b)
	case int64:
		i := parsedValue.(int64)
		field.SetInt(i)
	case string:
		s := parsedValue.(string)
		field.SetString(s)
	case float64:
		f := parsedValue.(float64)
		field.SetFloat(f)
	case uuid.UUID:
		u := parsedValue.(uuid.UUID)
		v := reflect.ValueOf(u)
		field.Set(v)
	}
	return nil
}

func JsonKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

type RelType int64

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Relationship struct {
	Id           uuid.UUID
	LeftBinding  RelType
	RightBinding RelType
	LeftModelId  uuid.UUID
	RightModelId uuid.UUID
	LeftName     string
	RightName    string
}

func (r Relationship) Left() Binding {
	return Binding{Relationship: r, Left: true}
}

func (r Relationship) Right() Binding {
	return Binding{Relationship: r, Left: false}
}

func JsonKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vId", strings.Title(strings.ToLower(key)))
}

type Model struct {
	Id                 uuid.UUID
	Name               string
	Attributes         map[string]Attribute
	LeftRelationships  []Relationship
	RightRelationships []Relationship
}

func (m Model) AttributeByName(name string) Attribute {
	a, ok := m.Attributes[name]
	if !ok {
		a, ok = SystemAttrs[name]
	}
	return a
}

func (m Model) GetBinding(name string) (Binding, error) {
	for _, b := range m.Bindings() {
		if b.Name() == name {
			return b, nil
		}
	}
	return Binding{}, ErrInvalidRelationship
}

func (m Model) Bindings() []Binding {
	var bs []Binding
	for _, r := range m.LeftRelationships {
		bs = append(bs, Binding{Relationship: r, Left: true})
	}
	for _, r := range m.RightRelationships {
		bs = append(bs, Binding{Relationship: r, Left: false})
	}
	return bs
}

type Binding struct {
	Relationship Relationship
	Left         bool
}

func (b Binding) HasField() bool {
	if b.Left {
		return (b.Relationship.LeftBinding == BelongsTo)
	} else {
		return (b.Relationship.RightBinding == BelongsTo)
	}
}

func (b Binding) Name() string {
	if b.Left {
		return b.Relationship.LeftName
	} else {
		return b.Relationship.RightName
	}
}

func (b Binding) ModelId() uuid.UUID {
	if b.Left {
		return b.Relationship.LeftModelId
	} else {
		return b.Relationship.RightModelId
	}
}

func (b Binding) Dual() Binding {
	return Binding{Relationship: b.Relationship, Left: !b.Left}
}

func (b Binding) RelType() RelType {
	if b.Left {
		return b.Relationship.LeftBinding
	} else {
		return b.Relationship.RightBinding
	}
}
