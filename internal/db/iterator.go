package db

import "errors"

type Iterator interface {
	Next() bool
	Value() interface{}
	Err() error
}

var Done = errors.New("done")
