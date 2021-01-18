package db

import (
	"fmt"
)

// RelLookupNode returns an iterator that returns the child rows for a given rel for a given parent
// it is the responsibility of the parent join node to set the results on the QR and apply aggregation
// the RelLookupNode does handle ordering
// the RelLookupNode does not handle filtering, so Aggregations may be applied
// this structure allows for offset/limit to be interposed between a join and a RelLookup
type RelLookupNode struct {
	tx          *holdTx
	rel         Relationship
	interfaceID ID
	iface       Interface

	order      []Sort
	projection Selection
}

func (rl *RelLookupNode) String() string {
	s := fmt.Sprintf("RelLookupNode{interface: %v, rel: %v, order: %v, projection: %v}",
		rl.iface.Name(), rl.rel.Name(), rl.order, rl.projection)
	return s
}

func (rl *RelLookupNode) Children() []Node {
	return []Node{}
}

func (rl *RelLookupNode) ResultIter(tx *holdTx, parentQR *QueryResult) (qrIterator, error) {
	rel, err := tx.Schema().GetRelationshipByID(rl.rel.ID())
	if err != nil {
		return nil, err
	}
	return &rlIterator{tx: rl.tx,
		parentQR:   parentQR,
		rel:        rel,
		order:      rl.order,
		projection: rl.projection,
	}, nil
}

type rlIterator struct {
	tx       *holdTx
	parentQR *QueryResult

	rel         Relationship
	interfaceID ID

	order      []Sort
	projection Selection

	pos    int
	values []Record
	err    error
}

func (i *rlIterator) Next() bool {
	if i.values == nil {
		if i.rel.Multi() {
			err := i.joinMany()
			if err != nil {
				i.err = err
				return false
			}
		} else {
			err := i.joinOne()
			if err != nil {
				i.err = err
				return false
			}
		}
	} else {
		// we've already loaded, so just return the next result
		i.pos++
	}

	if i.pos >= len(i.values) {
		i.err = Done
		return false
	}
	return true
}

func (i *rlIterator) Value() *QueryResult {
	rec := i.values[i.pos]
	qr := &QueryResult{Record: rec}
	project(qr, i.projection)
	return qr
}

func (i *rlIterator) Err() error {
	return i.err
}

func (i *rlIterator) joinMany() error {
	qr := i.parentQR
	relatedRecords, err := i.rel.LoadMany(qr.Record)
	if err != nil {
		return err
	}

	order(relatedRecords, i.order)

	i.values = relatedRecords
	return nil
}

func (i *rlIterator) joinOne() error {
	qr := i.parentQR
	relatedRecord, err := i.rel.LoadOne(qr.Record)
	if err != nil {
		if err == ErrNotFound {
			i.values = []Record{}
			return nil
		}
		return err
	}
	i.values = []Record{relatedRecord}
	return nil
}
