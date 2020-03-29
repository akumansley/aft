package db

import (
	"bytes"
	"fmt"
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
	idA := a.(Ider).GetId()
	idB := b.(Ider).GetId()
	result := bytes.Compare([]byte(idA), []byte(idB))

	return result
}

func (t *Table) Init() {
	t.t = gtreap.NewTreap(stringIdCompare)
}

func (t *Table) Upsert(item Ider) {
	t.t = t.t.Upsert(item, rand.Int()) // rand approximates balanced
}

func (t *Table) Get(id string) interface{} {
	it := t.t.Get(justId{id: id})
	return it
}

func (t *Table) printTree() {
	t.t.VisitAscend(t.t.Min(), func(i gtreap.Item) bool {
		fmt.Println("Visiting", i)
		return true
	})
}
