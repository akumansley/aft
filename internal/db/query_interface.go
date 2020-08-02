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
	"fmt"
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

func (qr *QueryResult) String() string {
	json, _ := qr.MarshalJSON()
	return string(json)
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

func (qr *QueryResult) GetChildRel(rel Relationship) []*QueryResult {
	if rel.Multi() {
		if v, ok := qr.ToMany[rel.Name()]; ok {
			return v
		}

	} else {
		if v, ok := qr.ToOne[rel.Name()]; ok {
			return []*QueryResult{v}
		}
	}
	return []*QueryResult{}
}

type ModelRef struct {
	interfaceID ID
	aliasID     uuid.UUID
	i           Interface
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

type JoinOperation struct {
	to ModelRef
	on RefRelationship
	jt joinType
}

func (j JoinOperation) String() string {
	return fmt.Sprintf("join: %v\t(%v)", j.to.i.Name(), j.to.aliasID)
}

func (q Q) String() string {
	var ifID string
	if q.main.Root == nil {
		ifID = "subquery"
	} else {
		ifID = q.main.Root.interfaceID.String()
	}
	return fmt.Sprintf(`Root: %v
Aggregations: %v
Joins: %v
Filters: %v
SetOps: %v`, ifID, q.main.Aggregations, q.main.Joins, q.main.Filters, q.main.SetOps)
}

type SetOpType int

const (
	or SetOpType = iota
	and
	not
)

type SetOperation struct {
	op       SetOpType
	branches []QBlock
}

func (j JoinOperation) IsToOne() bool {
	return !j.on.rel.Multi()
}

func (j JoinOperation) Key() string {
	return j.on.rel.Name()
}

type Q struct {
	tx   *holdTx
	main QBlock
}

func (tx *holdTx) Ref(interfaceID ID) ModelRef {
	i, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		panic("Bad ref")
	}
	return ModelRef{interfaceID, uuid.New(), i}
}

type QueryClause func(*QBlock)

func (tx *holdTx) Query(model ModelRef, clauses ...QueryClause) Q {
	qb := initQB()
	qb.Root = &model
	for _, c := range clauses {
		c(&qb)
	}
	return Q{tx: tx, main: qb}
}

func LeftJoin(to ModelRef, on RefRelationship) QueryClause {
	return func(qb *QBlock) {
		outer := on.from
		j := JoinOperation{to, on, leftJoin}
		joinList, ok := qb.Joins[outer.aliasID]
		if ok {
			qb.Joins[outer.aliasID] = append(joinList, j)
		} else {
			qb.Joins[outer.aliasID] = []JoinOperation{j}
		}
	}
}

func Join(to ModelRef, on RefRelationship) QueryClause {
	return func(qb *QBlock) {
		outer := on.from
		j := JoinOperation{to, on, innerJoin}
		joinList, ok := qb.Joins[outer.aliasID]
		if ok {
			qb.Joins[outer.aliasID] = append(joinList, j)
		} else {
			qb.Joins[outer.aliasID] = []JoinOperation{j}
		}

	}
}

func Filter(ref ModelRef, m Matcher) QueryClause {
	return func(qb *QBlock) {
		matcherList, ok := qb.Filters[ref.aliasID]
		if ok {
			qb.Filters[ref.aliasID] = append(matcherList, m)
		} else {
			qb.Filters[ref.aliasID] = []Matcher{m}
		}
	}
}

func Aggregate(ref ModelRef, a Aggregation) QueryClause {
	return func(qb *QBlock) {
		qb.Aggregations[ref.aliasID] = a
	}
}

func Or(ref ModelRef, branches ...QBlock) QueryClause {
	return SetOpClause(ref, or, branches...)
}

func Union(ref ModelRef, branches ...QBlock) QueryClause {
	return SetOpClause(ref, and, branches...)
}

func Not(ref ModelRef, branches ...QBlock) QueryClause {
	return SetOpClause(ref, not, branches...)
}

func SetOpClause(ref ModelRef, op SetOpType, branches ...QBlock) QueryClause {
	return func(qb *QBlock) {
		sos, ok := qb.SetOps[ref.aliasID]
		so := SetOperation{op, branches}
		if ok {
			qb.SetOps[ref.aliasID] = append(sos, so)
		} else {
			qb.SetOps[ref.aliasID] = []SetOperation{so}
		}
	}
}

func (q Q) All() []*QueryResult {
	results := q.main.runBlockRoot(q.tx)
	return results
}

func (q Q) One() (*QueryResult, error) {
	results := q.All()
	if len(results) == 0 {
		return nil, ErrNotFound
	}
	if len(results) != 1 {
		panic("Called one but got many")
	}
	return results[0], nil
}

func (q Q) OneRecord() (Record, error) {
	res, err := q.One()
	if err != nil {
		return nil, err
	}
	return res.Record, err
}

type QBlock struct {
	// TODO should we actually export these (vs exporting some setters)
	// null if this isn't a Root QB
	Root         *ModelRef
	Aggregations map[uuid.UUID]Aggregation
	Joins        map[uuid.UUID][]JoinOperation
	Filters      map[uuid.UUID][]Matcher
	SetOps       map[uuid.UUID][]SetOperation
}

func initQB() QBlock {
	return QBlock{
		Aggregations: map[uuid.UUID]Aggregation{},
		Filters:      map[uuid.UUID][]Matcher{},
		SetOps:       map[uuid.UUID][]SetOperation{},
		Joins:        map[uuid.UUID][]JoinOperation{}}
}

func Subquery(clauses ...QueryClause) QBlock {
	qb := initQB()
	for _, c := range clauses {
		c(&qb)
	}
	return qb
}

func NewBlock() QBlock {
	return initQB()
}
