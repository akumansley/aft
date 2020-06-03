package db

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool(t *testing.T) {
	var r interface{}
	var err error
	r, _ = boolValidatorFunc(false)
	assert.Equal(t, r, false)
	r, err = boolValidatorFunc("bad data")
	assert.Error(t, err)
	r, err = boolValidatorFunc(0)
	assert.Error(t, err)
}

func TestInt(t *testing.T) {
	var r interface{}
	var err error
	r, _ = intValidatorFunc(1)
	assert.Equal(t, r, int64(1))
	r, err = intValidatorFunc("bad data")
	assert.Error(t, err)
	r, err = intValidatorFunc(0.6)
	assert.Equal(t, r, int64(0))
	r, err = intValidatorFunc(int(6))
	assert.Equal(t, r, int64(6))
}

func TestEnum(t *testing.T) {
	var r interface{}
	var err error
	r, _ = enumValidatorFunc(1)
	assert.Equal(t, r, int64(1))
	r, err = enumValidatorFunc("bad data")
	assert.Error(t, err)
	r, err = enumValidatorFunc(0.6)
	assert.Equal(t, r, int64(0))
	r, err = enumValidatorFunc(int(6))
	assert.Equal(t, r, int64(6))
}

func TestString(t *testing.T) {
	var r interface{}
	var err error
	r, _ = stringValidatorFunc("this is a test")
	assert.Equal(t, r, "this is a test")
	r, err = stringValidatorFunc(5)
	assert.Error(t, err)
	r, _ = stringValidatorFunc("5")
	assert.Equal(t, r, "5")
}

func TestText(t *testing.T) {
	var r interface{}
	var err error
	r, _ = textValidatorFunc("this is a test")
	assert.Equal(t, r, "this is a test")
	r, err = textValidatorFunc(5)
	assert.Error(t, err)
	r, _ = textValidatorFunc("5")
	assert.Equal(t, r, "5")
}

func TestEmailAddress(t *testing.T) {
	var r interface{}
	var err error
	r, _ = emailAddressValidatorFunc("andrew.wansley@gmail.com")
	assert.Equal(t, r, "andrew.wansley@gmail.com")
	r, err = emailAddressValidatorFunc("wansley")
	assert.Error(t, err)
	r, err = emailAddressValidatorFunc("wansley@")
	assert.Error(t, err)
	r, err = emailAddressValidatorFunc("wansley@.c")
	assert.Error(t, err)
	r, err = emailAddressValidatorFunc("chase@a.co")
	assert.NoError(t, err)
}

func TestUUID(t *testing.T) {
	var r interface{}
	var err error

	r, _ = uuidValidatorFunc("d9c59e23-e050-4fc7-949d-8535ae8e3a49")
	u, _ := uuid.Parse("d9c59e23-e050-4fc7-949d-8535ae8e3a49")
	assert.Equal(t, r, u)
	r, err = uuidValidatorFunc("d9c59e23-e050-4fc7-949d-8535ae8e3a4")
	assert.Error(t, err)
	r, err = uuidValidatorFunc(123123)
	assert.Error(t, err)
	u, _ = uuid.Parse("39804353-56ee-463e-9b05-5e916ea293bd")
	r, _ = uuidValidatorFunc(u)
	assert.Equal(t, r, u)
}

func TestFloat(t *testing.T) {
	var r interface{}
	var err error
	r, _ = floatValidatorFunc(1)
	assert.Equal(t, r, float64(1))
	r, err = floatValidatorFunc("bad data")
	assert.Error(t, err)
	r, err = floatValidatorFunc(0.6)
	assert.Equal(t, r, 0.6)
	r, err = floatValidatorFunc(int(6))
	assert.Equal(t, r, float64(6))
}

func TestURL(t *testing.T) {
	var r interface{}
	var err error
	r, _ = URLValidatorFunc("https://www.google.com")
	assert.Equal(t, r, "https://www.google.com")
	r, _ = URLValidatorFunc("https://www.google.comdaddy")
	assert.Equal(t, r, "https://www.google.comdaddy")
	r, err = URLValidatorFunc("www.google.com")
	assert.Error(t, err)
	r, err = URLValidatorFunc("localhost:8080")
	assert.Error(t, err)
	r, err = URLValidatorFunc(5)
	assert.Error(t, err)
}
