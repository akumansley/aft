package db

import (
	"testing"

	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var datatypeTests = []struct {
	in          interface{}
	out         interface{}
	shouldError bool
	parser      func(context.Context, []interface{}) (interface{}, error)
}{
	{false, false, false, BoolFromJSON},
	{"bad data", false, true, BoolFromJSON},
	{0, false, true, BoolFromJSON},
	{1, int64(1), false, IntFromJSON},
	{"bad data", int64(1), true, IntFromJSON},
	{0.6, int64(0), false, IntFromJSON},
	{6, int64(6), false, IntFromJSON},
	{"test", "test", false, StringFromJSON},
	{5, "test", true, StringFromJSON},
	{"5", "5", false, StringFromJSON},
	{"d9c59e23-e050-4fc7-949d-8535ae8e3a49", uuid.MustParse("d9c59e23-e050-4fc7-949d-8535ae8e3a49"), false, UUIDFromJSON},
	{"d9c59e23-e050-4fc7-949d-8535ae8e3a4", "d9c59e23-e050-4fc7-949d-8535ae8e3a49", true, UUIDFromJSON},
	{321321, "d9c59e23-e050-4fc7-949d-8535ae8e3a49", true, UUIDFromJSON},
	{uuid.MustParse("d9c59e23-e050-4fc7-949d-8535ae8e3a49"), uuid.MustParse("d9c59e23-e050-4fc7-949d-8535ae8e3a49"), false, UUIDFromJSON},
	{1, float64(1), false, FloatFromJSON},
	{"bad data", float64(1), true, FloatFromJSON},
	{0.6, 0.6, false, FloatFromJSON},
	{int(6), float64(6), false, FloatFromJSON},
}

func TestParsers(t *testing.T) {
	for _, tt := range datatypeTests {
		r, err := tt.parser(context.Background(), []interface{}{tt.in})
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, r, tt.out)
		}
	}
}
