package db

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type TableAccessNode struct {
	interfaceID ID
	iface       Interface

	filters    []Matcher
	order      []Sort
	projection Selection
}

func (ta *TableAccessNode) String() string {
	s := fmt.Sprintf("TableAccessNode{interface: %v, filters: %v, order: %v, projection: %v}",
		ta.iface.Name(), ta.filters, ta.order, ta.projection)
	return s
}

func (ta *TableAccessNode) Children() []Node {
	return []Node{}
}

func (ta *TableAccessNode) ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error) {
	recs, err := tx.h.FindMany(ta.interfaceID, And(ta.filters...))
	if err != nil {
		return nil, err
	}
	order(recs, ta.order)

	return &taIterator{recs: recs, projection: ta.projection}, nil
}

func order(recs []Record, order []Sort) []Record {
	for _, step := range order {
		sort.Slice(recs, func(i, j int) bool {
			val := recs[i].MustGet(step.AttributeName)
			val2 := recs[j].MustGet(step.AttributeName)

			switch val.(type) {
			case int64:
				return val.(int64) < val2.(int64)
			case string:
				return val.(string) < val2.(string)
			case []byte:
				return bytes.Compare(val.([]byte), val2.([]byte)) < 0
			case float32:
				return val.(float32) < val2.(float32)
			case uuid.UUID:
				b1, _ := val.(uuid.UUID).MarshalBinary()
				b2, _ := val2.(uuid.UUID).MarshalBinary()
				return bytes.Compare(b1, b2) < 0
			case bool:
				panic("cannot order booleans")
			}
			panic("invalid type")
		})
	}
	return recs
}

type taIterator struct {
	recs       []Record
	projection Selection
	pos        int
	value      Record
	err        error
}

func (i *taIterator) Next() bool {
	if i.pos < len(i.recs) {
		i.value = i.recs[i.pos]
		i.pos++
		return true
	}
	i.err = Done
	return false
}

func project(qr *QueryResult, projection Selection) {
	if projection.selecting {
		qr.HideAll()
		for k, _ := range projection.fields {
			qr.Show(k)
		}
	}
}

func (i *taIterator) Value() *QueryResult {
	qr := &QueryResult{Record: i.value}
	project(qr, i.projection)
	return qr
}

func (i *taIterator) Err() error {
	return i.err
}
