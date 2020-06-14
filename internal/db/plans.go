package db

import (
	"github.com/google/uuid"
)

type framemaker struct {
	capacity int
	amap     map[uuid.UUID]int
}

func (fm framemaker) ix(aliasID uuid.UUID) int {
	return fm.amap[aliasID]
}

type intermediate interface {
	ID() uuid.UUID
}

type frame struct {
	entries []intermediate
}

func newFramemaker(q *Query) *framemaker {
	numScans := len(q.aset)
	fm := framemaker{capacity: numScans}
	ix := 0
	for k := range q.aset {
		fm.amap[k] = ix
		ix++
	}
	return &fm
}

type QueryResult struct {
	rows []frame
}

func (tx *holdTX) Execute(q *Query) (qr QueryResult, err error) {

	qr = QueryResult{rows: []frame{}}
	for it.Next() {
		val = it.Value()
		qr.rows = append(qr.rows, val)
	}
	err = it.Err()
	return
}

type plannode interface{}

type lookup struct {
	modelID uuid.UUID
	frameIx int
}

func plan(q *Query) PlanNode {
	fm := newFramemaker(q)
}

func planRelation(r relation) PlanNode {
	switch r.(type) {
	case table:
		t := r.(table)
	}
	return nil
}
