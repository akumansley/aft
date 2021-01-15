package db

import (
	"fmt"
)

type FilterNode struct {
	matcher Matcher
	inner   Node
}

func (c *FilterNode) String() string {
	return fmt.Sprintf("FilterNode{%v}", c.matcher)
}

func (c *FilterNode) Children() []Node {
	return []Node{c.inner}
}

func (c *FilterNode) ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error) {
	innerIter, err := c.inner.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}
	return &filterIterator{inner: innerIter, matcher: c.matcher, qr: qr}, nil
}

type filterIterator struct {
	qr      *QueryResult
	inner   qrIterator
	matcher Matcher
	value   *QueryResult
	err     error
}

func (i *filterIterator) Next() bool {
	for i.inner.Next() {
		value := i.inner.Value()
		match, err := i.matcher.Match(value.Record)
		if err != nil {
			i.err = err
			return false
		}
		if !match {
			continue
		}
		i.value = value
		return true
	}
	i.err = i.inner.Err()
	return false
}

func (i *filterIterator) Value() *QueryResult {
	return i.value
}

func (i *filterIterator) Err() error {
	return i.err
}
