package db

import (
	"github.com/google/uuid"
)

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

type framemaker struct {
	capacity int
	amap     map[uuid.UUID]int
}

type intermediate interface {
	ID() uuid.UUID
}

type frame struct {
	entries []intermediate
}

func (t *table) iter(tx Tx) (iterator, error) {
	recs := tx.FindMany(t.ref.modelID, And(s.where...))
	return &reciter{recs: recs, ix: 0}, nil
}

type group struct {
	key   uuid.UUID
	group []Record
}

func (g group) ID() uuid.UUID {
	return g.key
}

type groupiter struct {
	groups map[uuid.UUID][]Record
}

func (j *joinmany) iter(tx Tx) (iterator, error) {
	b := g.groupBy
	fk := b.Dual().Name()

	it, err := g.inner.iter(tx)
	if err != nil {
		return nil, err
	}
	hash := map[uuid.UUID][]Record{}
	for it.Next() {
		v := it.Value()
		k := v.GetFK(fk)
		ls, ok := hash[k]
		if ok {
			hash[k] = append(ls, v)
		} else {
			hash[k] = []Record{v}
		}
	}
	if it.Err() != nil {
		return nil, it.Err()
	}
	return &groupiter{groups: hash}, nil

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
	fm := newFramemaker(q)
	plan := q.r.plan(q)
	it := plan.iter()

	qr = QueryResult{rows: []frame{}}
	for it.Next() {
		val = it.Value()
		qr.rows = append(qr.rows, val)
	}
	err = it.Err()
	return
}
