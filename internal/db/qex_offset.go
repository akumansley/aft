package db

import "fmt"

type OffsetNode struct {
	offset int
	inner  Node
}

func (o *OffsetNode) String() string {
	return fmt.Sprintf("OffsetNode{offset: %v}", o.offset)
}

func (o *OffsetNode) Children() []Node {
	return []Node{o.inner}
}

func (o *OffsetNode) ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error) {
	innerIter, err := o.inner.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}
	return &offsetIterator{inner: innerIter, offset: o.offset}, nil
}

type offsetIterator struct {
	inner   qrIterator
	skipped int
	offset  int
}

func (i *offsetIterator) Next() bool {
	for i.skipped < i.offset {
		if i.inner.Next() {
			i.skipped++
		} else {
			return false
		}
	}

	i.inner.Next()
	return true
}

func (i *offsetIterator) Value() *QueryResult {
	return i.inner.Value()
}

func (i *offsetIterator) Err() error {
	return i.inner.Err()
}
