package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJavascript(t *testing.T) {
	result, _ := Javascript.FromJson(`print("my cool test\n\n"); 37;`)
	assert.Equal(t, 37.0, result)
}
