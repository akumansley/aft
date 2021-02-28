package db

import "fmt"

type JoinNode struct {
	aggregation Aggregation
	joinType    joinType
	rel         Relationship
	inner       Node
	outer       Node
	filters     []Matcher
}

func (j *JoinNode) String() string {
	s := fmt.Sprintf("JoinNode{%v on %v aggregation: %v}", j.joinType, j.rel.Name(), j.aggregation)
	return s
}

func (j *JoinNode) Children() []Node {
	return []Node{j.inner, j.outer}
}

func (j *JoinNode) ResultIter(tx *txWithContext, qr *QueryResult) (qrIterator, error) {
	innerIter, err := j.inner.ResultIter(tx, qr)
	if err != nil {
		return nil, err
	}

	return &joinIterator{tx: tx, inner: innerIter, outerNode: j.outer,
		aggregation: j.aggregation, joinType: j.joinType, rel: j.rel, filters: j.filters}, nil
}

type joinIterator struct {
	tx *txWithContext

	aggregation Aggregation
	joinType    joinType
	rel         Relationship
	inner       qrIterator
	outerNode   Node
	filters     []Matcher

	value *QueryResult
	pos   int
	err   error
}

func (i *joinIterator) loadInner() (qr *QueryResult, err error) {
	ok := i.inner.Next()
	if ok {
		qr = i.inner.Value()
	} else {
		err = i.inner.Err()
	}
	return
}

func (i *joinIterator) loadOuter(innerQR *QueryResult) (err error) {
	outer, err := i.outerNode.ResultIter(i.tx, innerQR)
	if err != nil {
		return
	}

	if i.rel.Multi() {
		return i.joinMany(innerQR, outer)
	} else {
		return i.joinOne(innerQR, outer)
	}
}

func (i *joinIterator) Next() bool {
	for {
		innerQR, err := i.loadInner()
		if err != nil {
			i.err = err
			return false
		}

		err = i.loadOuter(innerQR)
		if err != nil {
			i.err = err
			return false
		}
		if innerQR.isEmpty() {
			continue
		}
		i.value = innerQR

		return true
	}
}

func (i *joinIterator) Value() *QueryResult {
	return i.value
}

func (i *joinIterator) Err() error {
	return i.err
}

func (i *joinIterator) joinOne(innerQR *QueryResult, outer qrIterator) error {
	matcher := And(i.filters...)

	ok := outer.Next()
	if ok {
		outerQR := outer.Value()
		match, err := matcher.Match(outerQR.Record)
		if err != nil {
			return err
		}
		if match {
			innerQR.SetChildRelOne(i.rel.Name(), outerQR)
			return nil
		}
	}
	if outer.Err() != Done && outer.Err() != nil {
		return outer.Err()
	}

	// either the inner returned nothing
	// or the inner returned somenthing that didn't match
	if i.joinType == leftJoin {
		innerQR.SetChildRelOne(i.rel.Name(), &QueryResult{})
		return nil
	} else if i.joinType == innerJoin {
		innerQR.Empty()
		return nil
	} else {
		panic("invalid joinType")
	}
}

func (i *joinIterator) joinMany(innerQR *QueryResult, outer qrIterator) error {
	matcher := And(i.filters...)

	all := true
	matched := []*QueryResult{}
	for outer.Next() {
		outerQR := outer.Value()
		match, err := matcher.Match(outerQR.Record)
		if err != nil {
			return err
		}
		if match {
			matched = append(matched, outerQR)
		} else {
			all = false
		}
	}
	if outer.Err() != Done {
		return outer.Err()
	}

	if i.aggregation == Some {
		if len(matched) == 0 {
			innerQR.Empty()
		} else {
			innerQR.SetChildRelMany(i.rel.Name(), matched)
		}
	} else if i.aggregation == Every {
		if !all {
			innerQR.Empty()
		} else {
			innerQR.SetChildRelMany(i.rel.Name(), matched)
		}
	} else if i.aggregation == None {
		if len(matched) != 0 {
			innerQR.Empty()
		} else {
			innerQR.SetChildRelMany(i.rel.Name(), matched)
		}
	} else if i.aggregation == Include {
		innerQR.SetChildRelMany(i.rel.Name(), matched)
	} else {
		panic("invalid aggregation")
	}

	return nil
}
