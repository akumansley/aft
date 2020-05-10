package er

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"errors"
	"fmt"
	"github.com/hashicorp/go-immutable-radix"
)

var (
	ErrHold     = errors.New("hold-error")
	ErrNotFound = fmt.Errorf("%w: not found", ErrHold)
)

type Hold struct {
	t *iradix.Tree
}

func New() *Hold {
	return &Hold{t: iradix.New()}
}

func (h *Hold) FindOne(table string, q q.Matcher) (model.Record, error) {
	it := h.t.Root().Iterator()
	it.SeekPrefix([]byte(table))

	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		rec := val.(model.Record)
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

type Iterator interface {
	Next() (model.Record, bool)
}

type MatchIter struct {
	q  q.Matcher
	it *iradix.Iterator
}

func (mi MatchIter) Next() (model.Record, bool) {
	for _, val, ok := mi.it.Next(); ok; _, val, ok = mi.it.Next() {
		rec := val.(model.Record)
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

func (h *Hold) IterMatches(table string, q q.Matcher) Iterator {
	it := h.t.Root().Iterator()
	it.SeekPrefix([]byte(table))
	return MatchIter{q: q, it: it}
}

func makeKey(rec model.Record) []byte {
	ub, _ := rec.Id().MarshalBinary()
	bytes := append(append([]byte(rec.Type()), []byte("/")...), ub...)
	return bytes
}

func (h *Hold) Insert(rec model.Record) *Hold {
	newTree, _, _ := h.t.Insert(makeKey(rec), rec)
	return &Hold{t: newTree}

}

func (h *Hold) PrintTree() {
	fmt.Printf("print tree:\n")
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		fmt.Printf("%v:%v\n", string(k), v)
	}
	fmt.Printf("done printing\n")
}
