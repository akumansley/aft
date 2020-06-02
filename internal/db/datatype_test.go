package db

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	//	"fmt"
)

func TestBool(t *testing.T) {
	var r interface{}
	var err error
	r, _ = boolFromJSONFunc(false)
	assert.Equal(t, r, false)
	r, err = boolFromJSONFunc("bad data")
	assert.Error(t, err)
	r, err = boolFromJSONFunc(0)
	assert.Error(t, err)
}

func TestInt(t *testing.T) {
	var r interface{}
	var err error
	r, _ = intFromJSONFunc(1)
	assert.Equal(t, r, int64(1))
	r, err = intFromJSONFunc("bad data")
	assert.Error(t, err)
	r, err = intFromJSONFunc(0.6)
	assert.Equal(t, r, int64(0))
	r, err = intFromJSONFunc(int(6))
	assert.Equal(t, r, int64(6))
}

func TestEnum(t *testing.T) {
	var r interface{}
	var err error
	r, _ = enumFromJSONFunc(1)
	assert.Equal(t, r, int64(1))
	r, err = enumFromJSONFunc("bad data")
	assert.Error(t, err)
	r, err = enumFromJSONFunc(0.6)
	assert.Equal(t, r, int64(0))
	r, err = enumFromJSONFunc(int(6))
	assert.Equal(t, r, int64(6))
}

func TestString(t *testing.T) {
	var r interface{}
	var err error
	r, _ = stringFromJSONFunc("this is a test")
	assert.Equal(t, r, "this is a test")
	r, err = stringFromJSONFunc(5)
	assert.Error(t, err)
	r, _ = stringFromJSONFunc("5")
	assert.Equal(t, r, "5")
}

func TestText(t *testing.T) {
	var r interface{}
	var err error
	r, _ = textFromJSONFunc("this is a test")
	assert.Equal(t, r, "this is a test")
	r, err = textFromJSONFunc(5)
	assert.Error(t, err)
	r, _ = textFromJSONFunc("5")
	assert.Equal(t, r, "5")
}

func TestEmailAddress(t *testing.T) {
	var r interface{}
	var err error
	r, _ = emailAddressFromJSONFunc("andrew.wansley@gmail.com")
	assert.Equal(t, r, "andrew.wansley@gmail.com")
	r, err = emailAddressFromJSONFunc("wansley")
	assert.Error(t, err)
	r, err = emailAddressFromJSONFunc("wansley@")
	assert.Error(t, err)
	r, err = emailAddressFromJSONFunc("wansley@.c")
	assert.Error(t, err)
	r, err = emailAddressFromJSONFunc("chase@a.co")
	assert.NoError(t, err)
}

func TestUUID(t *testing.T) {
	var r interface{}
	var err error

	r, _ = uuidFromJSONFunc("d9c59e23-e050-4fc7-949d-8535ae8e3a49")
	u, _ := uuid.Parse("d9c59e23-e050-4fc7-949d-8535ae8e3a49")
	assert.Equal(t, r, u)
	r, err = uuidFromJSONFunc("d9c59e23-e050-4fc7-949d-8535ae8e3a4")
	assert.Error(t, err)
	r, err = uuidFromJSONFunc(123123)
	assert.Error(t, err)
	u, _ = uuid.Parse("39804353-56ee-463e-9b05-5e916ea293bd")
	r, _ = uuidFromJSONFunc(u)
	assert.Equal(t, r, u)
}

func TestFloat(t *testing.T) {
	var r interface{}
	var err error
	r, _ = floatFromJSONFunc(1)
	assert.Equal(t, r, float64(1))
	r, err = floatFromJSONFunc("bad data")
	assert.Error(t, err)
	r, err = floatFromJSONFunc(0.6)
	assert.Equal(t, r, 0.6)
	r, err = floatFromJSONFunc(int(6))
	assert.Equal(t, r, float64(6))
}

func TestURL(t *testing.T) {
	var r interface{}
	var err error
	r, _ = URLFromJSONFunc("https://www.google.com")
	assert.Equal(t, r, "https://www.google.com")
	r, _ = URLFromJSONFunc("https://www.google.comdaddy")
	assert.Equal(t, r, "https://www.google.comdaddy")
	r, err = URLFromJSONFunc("www.google.com")
	assert.Error(t, err)
	r, err = URLFromJSONFunc("localhost:8080")
	assert.Error(t, err)
	r, err = URLFromJSONFunc(5)
	assert.Error(t, err)
}
