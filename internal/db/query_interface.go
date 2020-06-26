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
	modelID ModelID
	aliasID uuid.UUID
	model   Model
}

type RefRelationship struct {
	from ModelRef
	rel  Relationship
}

func (ref ModelRef) Rel(rel Relationship) RefRelationship {
	return RefRelationship{ref, rel}
}

type joinType int

const (
	leftJoin joinType = iota
	innerJoin
)

type join struct {
	to ModelRef
	on RefRelationship
	jt joinType
}

type setoperation int

const (
	or setoperation = iota
	and
	not
)

type setop struct {
	op       setoperation
	branches []QBlock
}

func (j join) IsToOne() bool {
	return !j.on.rel.Multi
}

func (j join) Key() string {
	return j.on.rel.Name
}

type Q struct {
	tx   *holdTx
	main QBlock
}

func (tx *holdTx) Ref(modelID ModelID) ModelRef {
	model, _ := tx.GetModelByID(modelID)
	return ModelRef{modelID, uuid.New(), model}
}

func (tx *holdTx) Query(model ModelRef) Q {
	qb := initQB()
	qb.root = &model
	return Q{tx: tx, main: qb}
}

func (q Q) AsBlock() QBlock {
	return q.main
}

func (q Q) SetMainBlock(qb QBlock) {
	q.main = qb
}

func (q Q) LeftJoin(to ModelRef, on RefRelationship) Q {
	q.main = q.main.LeftJoin(to, on)
	return q
}

func (q Q) Join(to ModelRef, on RefRelationship) Q {
	q.main = q.main.Join(to, on)
	return q
}

func (q Q) Filter(ref ModelRef, m Matcher) Q {
	q.main = q.main.Filter(ref, m)
	return q
}

func (q Q) Aggregate(ref ModelRef, a Aggregation) Q {
	q.main = q.main.Aggregate(ref, a)
	return q
}

func (q Q) Or(ref ModelRef, branches ...QBlock) Q {
	q.main = q.main.Or(ref, branches...)
	return q
}

func (q Q) And(ref ModelRef, branches ...QBlock) Q {
	q.main = q.main.And(ref, branches...)
	return q
}

func (q Q) Not(ref ModelRef, branches ...QBlock) Q {
	q.main = q.main.Not(ref, branches...)
	return q
}

func (q Q) All() []*QueryResult {
	results := q.main.runBlockRoot(q.tx)
	return results
}

type QBlock struct {
	// null if this isn't a root QB
	root         *ModelRef
	aggregations map[uuid.UUID]Aggregation
	joins        map[uuid.UUID][]join
	sargs        map[uuid.UUID][]Matcher
	setops       map[uuid.UUID][]setop
}

func initQB() QBlock {
	return QBlock{
		sargs:        map[uuid.UUID][]Matcher{},
		aggregations: map[uuid.UUID]Aggregation{},
		setops:       map[uuid.UUID][]setop{},
		joins:        map[uuid.UUID][]join{}}
}

func NewBlock() QBlock {
	return initQB()
}

func Filter(ref ModelRef, m Matcher) QBlock {
	qb := initQB()
	qb = qb.Filter(ref, m)
	return qb
}

func (qb QBlock) Filter(ref ModelRef, m Matcher) QBlock {
	matcherList, ok := qb.sargs[ref.aliasID]
	if ok {
		qb.sargs[ref.aliasID] = append(matcherList, m)
	} else {
		qb.sargs[ref.aliasID] = []Matcher{m}
	}
	return qb
}

func Join(to ModelRef, on RefRelationship) QBlock {
	qb := initQB()
	qb = qb.Join(to, on)
	return qb
}

func (qb QBlock) LeftJoin(to ModelRef, on RefRelationship) QBlock {
	outer := on.from
	j := join{to, on, leftJoin}
	joinList, ok := qb.joins[outer.aliasID]
	if ok {
		qb.joins[outer.aliasID] = append(joinList, j)
	} else {
		qb.joins[outer.aliasID] = []join{j}
	}
	return qb
}

func (qb QBlock) Join(to ModelRef, on RefRelationship) QBlock {
	outer := on.from
	j := join{to, on, innerJoin}
	joinList, ok := qb.joins[outer.aliasID]
	if ok {
		qb.joins[outer.aliasID] = append(joinList, j)
	} else {
		qb.joins[outer.aliasID] = []join{j}
	}
	return qb
}

func Aggregate(ref ModelRef, a Aggregation) QBlock {
	qb := initQB()
	qb = qb.Aggregate(ref, a)
	return qb
}

func (qb QBlock) Aggregate(ref ModelRef, a Aggregation) QBlock {
	qb.aggregations[ref.aliasID] = a
	return qb
}

func (qb QBlock) setOp(ref ModelRef, op setoperation, branches ...QBlock) QBlock {
	sos, ok := qb.setops[ref.aliasID]
	so := setop{op, branches}
	if ok {
		qb.setops[ref.aliasID] = append(sos, so)
	} else {
		qb.setops[ref.aliasID] = []setop{so}
	}
	return qb
}

func (qb QBlock) Or(ref ModelRef, branches ...QBlock) QBlock {
	return qb.setOp(ref, or, branches...)
}

func (qb QBlock) And(ref ModelRef, branches ...QBlock) QBlock {
	return qb.setOp(ref, and, branches...)
}

func (qb QBlock) Not(ref ModelRef, branches ...QBlock) QBlock {
	return qb.setOp(ref, not, branches...)
}
