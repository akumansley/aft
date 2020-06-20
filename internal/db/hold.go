package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-immutable-radix"
)

var (
	ErrNotFound = fmt.Errorf("%w: not found", ErrData)
)

type Hold struct {
	t *iradix.Tree
}

func NewHold() *Hold {
	return &Hold{t: iradix.New()}
}

func (h *Hold) FindOne(modelID uuid.UUID, q Matcher) (Record, error) {
	mb, _ := modelID.MarshalBinary()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)

	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return nil, err
		}
		if match {
			return rec, nil
		}
	}
	return nil, ErrNotFound
}

type MatchIter struct {
	q  Matcher
	it *iradix.Iterator
}

func (mi MatchIter) Next() (Record, bool) {
	for _, val, ok := mi.it.Next(); ok; _, val, ok = mi.it.Next() {
		rec := val.(Record)
		match, err := mi.q.Match(rec)
		if err != nil {
			return nil, false
		}
		if match {
			return rec, true
		}
	}
	return nil, false
}

func (h *Hold) IterMatches(modelID uuid.UUID, q Matcher) Iterator {
	mb, _ := modelID.MarshalBinary()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)
	return MatchIter{q: q, it: it}
}

func (h *Hold) FindMany(modelID uuid.UUID, q Matcher) ([]Record, error) {
	mb, _ := modelID.MarshalBinary()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)
	hits := []Record{}
	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return hits, err
		}
		if match {
			hits = append(hits, rec)
		}
	}
	return hits, nil
}

func makeKey(rec Record) []byte {
	rb, _ := rec.ID().MarshalBinary()
	mb, _ := rec.Model().ID.MarshalBinary()

	bytes := append(append(mb, []byte("/")...), rb...)
	return bytes
}

func (h *Hold) Insert(rec Record) *Hold {
	newTree, _, _ := h.t.Insert(makeKey(rec), rec)
	return &Hold{t: newTree}
}

func (h *Hold) Delete(rec Record) *Hold {
	newTree, _, _ := h.t.Delete(makeKey(rec))
	return &Hold{t: newTree}
}

type RootIter struct {
	it *iradix.Iterator
}

func (ri RootIter) Next() (Record, bool) {
	for _, val, ok := ri.it.Next(); ok; _, val, ok = ri.it.Next() {
		rec := val.(Record)
		return rec, true
	}
	return nil, false
}

func (h *Hold) Iterator() Iterator {
	it := h.t.Root().Iterator()
	return RootIter{it: it}
}

func (h *Hold) PrintTree() {
	fmt.Printf("print tree:\n")
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		fmt.Printf("%v:%v\n", string(k), v)
	}
	fmt.Printf("done printing\n")
}
