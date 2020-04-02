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

type Table interface {
	Init()
	Put(Ider)
	Get(string) interface{}
	Query(string) []interface{}
}

type TreapTable struct {
	t *gtreap.Treap
	i Index // TODO factor this out
}

func stringIdCompare(a, b interface{}) int {
	idA := a.(Ider).GetId()
	idB := b.(Ider).GetId()
	result := bytes.Compare([]byte(idA), []byte(idB))

	return result
}

func (t *TreapTable) Init() {
	t.t = gtreap.NewTreap(stringIdCompare)
	t.i = &BleveIndex{}
	t.i.Init()
}

func (t *TreapTable) Put(item Ider) {
	t.t = t.t.Upsert(item, rand.Int()) // rand approximates balanced
	t.i.Index(item)
}

func (t *TreapTable) Get(id string) interface{} {
	it := t.t.Get(justId{id: id})
	return it
}

func (t *TreapTable) Query(q string) []interface{} {
	searchResults := t.i.Query(q)
	results := make([]interface{}, searchResults.Total)
	for i, hit := range searchResults.Hits {
		results[i] = t.Get(hit.ID)
	}
	return results
}

func (t *TreapTable) printTree() {
	t.t.VisitAscend(t.t.Min(), func(i gtreap.Item) bool {
		fmt.Println("Visiting", i)
		return true
	})
}
