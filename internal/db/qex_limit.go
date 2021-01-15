package db

import "fmt"

type LimitNode struct {
	limit int
	inner Node
}

func (l *LimitNode) String() string {
	return fmt.Sprintf("LimitNode{limit: %v}", l.limit)
}

func (l *LimitNode) Children() []Node {
	return []Node{l.inner}
}

func (l *LimitNode) ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error) {
	innerIter, err := l.inner.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}
	return &limitIterator{inner: innerIter, limit: l.limit}, nil
}

type limitIterator struct {
	inner qrIterator
	count int
	limit int
	err   error
}

func (i *limitIterator) Next() bool {
	if i.count == i.limit {
		i.err = Done
		return false
	}
	i.inner.Next()
	i.count++
	return true
}

func (i *limitIterator) Value() *QueryResult {
	return i.inner.Value()
}

func (i *limitIterator) Err() error {
	if i.err != nil {
		return i.err
	}
	return i.inner.Err()
}
