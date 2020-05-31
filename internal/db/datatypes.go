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

type Datatype struct {
	Id       uuid.UUID
	FromJson func(interface{}) (interface{}, error)
	Type     func() interface{}
}

var Bool = Datatype{
	Id:       uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	FromJson: boolFromJson,
	Type:     boolType,
}
var Int = Datatype{
	Id:       uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	FromJson: intFromJson,
	Type:     intType,
}
var Enum = Datatype{
	Id:       uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	FromJson: enumFromJson,
	Type:     enumType,
}
var String = Datatype{
	Id:       uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	FromJson: stringFromJson,
	Type:     stringType,
}
var Text = Datatype{
	Id:       uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	FromJson: textFromJson,
	Type:     textType,
}
var EmailAddress = Datatype{
	Id:       uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	FromJson: emailAddressFromJson,
	Type:     emailAddressType,
}
var UUID = Datatype{
	Id:       uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	FromJson: uuidFromJson,
	Type:     uuidType,
}
var Float = Datatype{
	Id:       uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	FromJson: floatFromJson,
	Type:     floatType,
}
var URL = Datatype{
	Id:       uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	FromJson: urlFromJson,
	Type:     urlType,
}

var datatypeMap map[uuid.UUID]Datatype = map[uuid.UUID]Datatype{
	Bool.Id:         Bool,
	Int.Id:          Int,
	Enum.Id:         Enum,
	String.Id:       String,
	Text.Id:         Text,
	EmailAddress.Id: EmailAddress,
	UUID.Id:         UUID,
	Float.Id:        Float,
	URL.Id:          URL,
}

//booleans
func boolFromJson(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

func boolType() interface{} {
	return false
}

//Integers
func intFromJson(value interface{}) (interface{}, error) {
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

func intType() interface{} {
	return int64(0)
}

//Enum
func enumFromJson(value interface{}) (interface{}, error) {
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

func enumType() interface{} {
	return int64(0)
}

//String
func stringFromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func stringType() interface{} {
	return ""
}

//Text
func textFromJson(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected text got %T", ErrValue, value)
	}
	return s, nil
}

func textType() interface{} {
	return ""
}

//EmailAddress
//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func emailAddressFromJson(value interface{}) (interface{}, error) {
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

func emailAddressType() interface{} {
	return ""
}

//UUID
func uuidFromJson(value interface{}) (interface{}, error) {
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

func uuidType() interface{} {
	return uuid.UUID{}
}

//Float
func floatFromJson(value interface{}) (interface{}, error) {
	f, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
	}
	return f, nil
}

func floatType() interface{} {
	return 0.0
}

//URL
func urlFromJson(value interface{}) (interface{}, error) {
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

func urlType() interface{} {
	return ""
}
