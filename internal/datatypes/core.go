package datatypes

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"regexp"
)

var (
	ErrValue = fmt.Errorf("invalid value for type")
)

var FunctionMap map[string]func(interface{}) (interface{}, error) = map[string]func(interface{}) (interface{}, error){
	"boolValidator":         boolValidatorFunc,
	"intValidator":          intValidatorFunc,
	"enumValidator":         enumValidatorFunc,
	"stringValidator":       stringValidatorFunc,
	"textValidator":         textValidatorFunc,
	"emailAddressValidator": emailAddressValidatorFunc,
	"uuidValidator":         uuidValidatorFunc,
	"floatValidator":        floatValidatorFunc,
	"URLValidator":          URLValidatorFunc,
}

func boolValidatorFunc(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

func intValidatorFunc(value interface{}) (interface{}, error) {
	return intEnumValidatorFunc(value, "int")
}

func enumValidatorFunc(value interface{}) (interface{}, error) {
	return intEnumValidatorFunc(value, "enum")
}

func intEnumValidatorFunc(value interface{}, t string) (interface{}, error) {
	switch value.(type) {
	case float64:
		return int64(value.(float64)), nil
	case int:
		return int64(value.(int)), nil
	case int64:
		return value, nil
	}
	return nil, fmt.Errorf("%w: expected %s got %T", ErrValue, t, value)

}

func stringValidatorFunc(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

func textValidatorFunc(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected text got %T", ErrValue, value)
	}
	return s, nil
}

// Email Address datatype uses following regex to validate emails.
//https://www.alexedwards.net/blog/validation-snippets-for-go#email-validation
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func matchEmail(s string) bool {
	return rxEmail.MatchString(s)
}

func emailAddressValidatorFunc(value interface{}) (interface{}, error) {
	es, ok := value.(string)
	if ok {
		if (len(es) > 254 || !matchEmail(es)) && len(es) != 0 {
			return nil, fmt.Errorf("%w: expected email address got %v", ErrValue, es)
		}
	} else {
		return nil, fmt.Errorf("%w: expected email address got %T", ErrValue, value)
	}
	return es, nil
}

func uuidValidatorFunc(value interface{}) (interface{}, error) {
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

func floatValidatorFunc(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int64:
		return float64(value.(int64)), nil
	case int:
		return float64(value.(int)), nil
	case float64:
		return value, nil
	}
	return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
}

func URLValidatorFunc(value interface{}) (interface{}, error) {
	us, ok := value.(string)
	if ok {
		u, err := url.Parse(us)
		if err != nil {
			return nil, fmt.Errorf("%w: expected URL got %s", ErrValue, u)
		} else if u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("%w: expected URL got %s", ErrValue, u)
		} else if u.Scheme != "http" && u.Scheme != "https" {
			return nil, fmt.Errorf("%w: expected URL got %s", ErrValue, u)
		}
	} else {
		return nil, fmt.Errorf("%w: expected URL got %T", ErrValue, value)
	}
	return us, nil
}
