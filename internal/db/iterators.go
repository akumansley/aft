package db

import ()

type iterator interface {
	Next() bool
	Value() interface{}
	Err() error
}

type frameiter struct {
	frames []frame
	ix     int
	value  frame
	err    error
}

func (i *frameiter) Value() interface{} {
	return i.value
}
func (i *frameiter) Err() error {
	return i.err
}

func (i *frameiter) Next() bool {
	if i.ix < len(i.frames) {
		i.ix++
		i.value = i.frames[i.ix-1]
		return true
	}
	return false
}

type reciter struct {
	recs  []Record
	ix    int
	value Record
	err   error
}

func (i *reciter) Value() interface{} {
	return i.value
}
func (i *reciter) Err() error {
	return i.err
}

func (i *reciter) Next() bool {
	if i.ix < len(i.recs) {
		i.ix++
		i.value = i.recs[i.ix-1]
		return true
	}
	return false
}
