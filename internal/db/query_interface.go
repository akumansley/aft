// basic operation:
// user = db.Ref(modelID)
// db.Query(user).Join(post, user.Rel["posts"]).Filter(post, db.Eq("foo", "bar")).All()
//
// or:
//
// q1 := db.Query(user).Filter(user, db.Eq("age", 32)).Or([
// 	db.Filter(user, db.Eq("name", "Andrew")),
// 	db.Filter(user, db.Eq("name", "Chase")).Join(posts, user.Rel("posts")).Filter(post, db.Eq("text", "hello")),
// ])
//
// db.Query(user).Join(post, user.Rel("posts")).Filter(post, db.Eq("text", "hello")).All()
// db.Query(user).Join(post, user.Rel("posts")).Filter(post, db.Eq("text", "goodbye")).All()

package db

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Aggregation int

const (
	Every Aggregation = iota
	Some
	None
	Include
)

type QueryResult struct {
	Record Record
	ToOne  map[string]*QueryResult
	ToMany map[string][]*QueryResult
}

func (qr *QueryResult) isEmpty() bool {
	return qr.Record == nil
}

func (qr *QueryResult) Empty() {
	qr.Record = nil
}

func (qr *QueryResult) MarshalJSON() ([]byte, error) {
	if qr.Record == nil {
		return json.Marshal(nil)
	}
	data := qr.Record.Map()
	for k, v := range qr.ToOne {
		data[k] = v
	}
	for k, v := range qr.ToMany {
		data[k] = v
	}
	return json.Marshal(data)
}

type ModelRef struct {
	modelID uuid.UUID
	aliasID uuid.UUID
	model   Model
}

type RefBinding struct {
	from ModelRef
	b    Binding
}

func (ref ModelRef) Rel(name string) RefBinding {
	b, _ := ref.model.GetBinding(name)
	return RefBinding{ref, b}
}

type joinType int

const (
	leftJoin joinType = iota
	innerJoin
)

type join struct {
	to ModelRef
	on RefBinding
	jt joinType
}

func (j join) Key() string {
	return j.on.b.Name()
}

func (j join) IsToOne() bool {
	return j.on.b.RelType() == HasOne || j.on.b.RelType() == BelongsTo
}

type Q struct {
	root         ModelRef
	aggregations map[uuid.UUID]Aggregation
	joins        map[uuid.UUID][]join
	tx           *holdTx
	// "search args" after system R
	sargs map[uuid.UUID][]Matcher
}

func (tx *holdTx) Ref(modelID uuid.UUID) ModelRef {
	model, _ := tx.GetModelByID(modelID)
	return ModelRef{modelID, uuid.New(), model}
}

func (tx *holdTx) Query(model ModelRef) Q {
	return Q{root: model, tx: tx,
		sargs:        map[uuid.UUID][]Matcher{},
		aggregations: map[uuid.UUID]Aggregation{},
		joins:        map[uuid.UUID][]join{}}
}

func (q Q) Join(to ModelRef, on RefBinding) Q {
	outer := on.from
	j := join{to, on, innerJoin}
	joinList, ok := q.joins[outer.aliasID]
	if ok {
		q.joins[outer.aliasID] = append(joinList, j)
	} else {
		q.joins[outer.aliasID] = []join{j}
	}
	return q
}

func (q Q) Filter(ref ModelRef, m Matcher) Q {
	matcherList, ok := q.sargs[ref.aliasID]
	if ok {
		q.sargs[ref.aliasID] = append(matcherList, m)
	} else {
		q.sargs[ref.aliasID] = []Matcher{m}
	}
	return q
}

func (q Q) Aggregate(ref ModelRef, a Aggregation) Q {
	q.aggregations[ref.aliasID] = a
	return q
}

func (q Q) All() []*QueryResult {
	var results []*QueryResult
	matchers := q.sargs[q.root.aliasID]
	results = q.performScan(q.root.modelID, And(matchers...))
	results = q.performJoins(results, q.root.aliasID)
	results = filterEmpty(results)
	return results
}
