package db

import (
	"bytes"
	"github.com/steveyen/gtreap"
	"math/rand"
)

type Ider interface {
	GetId() string
}
type justId struct {
	id string
}

func (i justId) GetId() string {
	return i.id
}

type Table struct {
	t *gtreap.Treap
}

func stringIdCompare(a, b interface{}) int {
	return bytes.Compare([]byte(a.(Ider).GetId()), []byte(b.(Ider).GetId()))
}

func (t *Table) Init() {
	t.t = gtreap.NewTreap(stringIdCompare)
}

func (t *Table) Upsert(id string, item Ider) {
	t.t = t.t.Upsert(item, rand.Int()) // rand approximates balanced
}

func (t *Table) Get(id string) interface{} {
	it := t.t.Get(justId{id: id})
	return it
}
