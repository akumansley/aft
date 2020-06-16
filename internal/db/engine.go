package db

import (
	"encoding/json"
	"github.com/google/uuid"
)

type QueryResult struct {
	Record        Record
	SingleRelated map[string]QueryResult
	MultiRelated  map[string][]QueryResult
}

func (ir QueryResult) MarshalJSON() ([]byte, error) {
	data := ir.Record.Map()
	for k, v := range ir.SingleRelated {
		data[k] = v
	}
	for k, v := range ir.MultiRelated {
		data[k] = v
	}
	return json.Marshal(data)
}

type group struct {
	key   uuid.UUID
	group []Record
}

func (g group) ID() uuid.UUID {
	return g.key
}

func (tx *holdTx) Execute(q *Query) (qr QueryResult, err error) {
	pn := plan(q)
	it := plan.iter(tx)

	for it.Next() {
		val = it.Value()
		qr.rows = append(qr.rows, val)
	}
	err = it.Err()
	return
}
