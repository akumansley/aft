package datatypes

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var datatypeTests = []struct {
		in interface{}
		out interface{}
		shouldError bool
		parser func(interface{}) (interface{}, error)
}{
	{false, false, false, boolValidatorFunc},
	{"bad data", false, true, boolValidatorFunc},
	{0, false, true, boolValidatorFunc},
	{1, int64(1), false, intValidatorFunc},
	{"bad data", int64(1), true, intValidatorFunc},
	{0.6, int64(0), false, intValidatorFunc},
	{6, int64(6), false, intValidatorFunc},
	{1, int64(1), false, enumValidatorFunc},
	{"bad data", int64(1), true, enumValidatorFunc},
	{0.6, int64(0), false, enumValidatorFunc},
	{6, int64(6), false, enumValidatorFunc},
	{"test", "test", false, stringValidatorFunc},
	{5, "test", true, stringValidatorFunc},
	{"5", "5", false, stringValidatorFunc},
	{"test", "test", false, textValidatorFunc},
	{5, "test", true, textValidatorFunc},
	{"5", "5", false, textValidatorFunc},
	{"andrew.wansley@gmail.com", "andrew.wansley@gmail.com", false, emailAddressValidatorFunc},
	{"wansley", "andrew.wansley@gmail.com", true, emailAddressValidatorFunc},
	{"wansley@", "andrew.wansley@gmail.com", true, emailAddressValidatorFunc},
	{"wansley@.c", "andrew.wansley@gmail.com", true, emailAddressValidatorFunc},
	{"chase@a.c", "chase@a.c", false, emailAddressValidatorFunc},
	{"d9c59e23-e050-4fc7-949d-8535ae8e3a49",getFirst(uuid.Parse("d9c59e23-e050-4fc7-949d-8535ae8e3a49")),false, uuidValidatorFunc},
	{"d9c59e23-e050-4fc7-949d-8535ae8e3a4","d9c59e23-e050-4fc7-949d-8535ae8e3a49",true,uuidValidatorFunc},
	{321321,"d9c59e23-e050-4fc7-949d-8535ae8e3a49",true,uuidValidatorFunc},
	{getFirst(uuid.Parse("d9c59e23-e050-4fc7-949d-8535ae8e3a49")),getFirst(uuid.Parse("d9c59e23-e050-4fc7-949d-8535ae8e3a49")),false, uuidValidatorFunc},
	{1, float64(1), false, floatValidatorFunc},
	{"bad data", float64(1), true, floatValidatorFunc},
	{0.6, 0.6, false, floatValidatorFunc},
	{int(6), float64(6), false, floatValidatorFunc},
	{"https://www.google.com","https://www.google.com", false, URLValidatorFunc},
	{"https://www.google.comdaddy","https://www.google.comdaddy", false, URLValidatorFunc},
	{"www.google.com","www.google.com", true, URLValidatorFunc},
	{"localhost:8080","www.google.com", true, URLValidatorFunc},
	{5,"www.google.com", true, URLValidatorFunc},
}

func getFirst(a, b interface{}) interface{}{
	return a
}

func TestParsers(t *testing.T) {
	for _, tt := range datatypeTests{
		r, err := tt.parser(tt.in)
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, r, tt.out)
		}
	}
}
