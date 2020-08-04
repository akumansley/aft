// basic operation:
// user = db.Ref(modelID)
// tx.Query(user).Join(post, user.Rel["posts"]).Filter(post, db.Eq("foo", "bar")).All()
//
// or:
//
// q1 := tx.Query(user).Filter(user, db.Eq("age", 32)).Or([
// 	db.Filter(user, db.Eq("name", "Andrew")),
// 	db.Filter(user, db.Eq("name", "Chase")).Join(posts, user.Rel("posts")).Filter(post, db.Eq("text", "hello")),
// ])
//
// tx.Query(user).Join(post, user.Rel("posts")).Filter(post, db.Eq("text", "hello")).All()
// tx.Query(user).Join(post, user.Rel("posts")).Filter(post, db.Eq("text", "goodbye")).All()

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

func (a Aggregation) String() string {
	switch a {
	case Every:
		return "Every"
	case Some:
		return "Some"
	case None:
		return "None"
	case Include:
		return "Include"
	default:
		panic("Invalid aggregation")
	}
}

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
	data["type"] = qr.Record.Type()
	return json.Marshal(data)
}

func (qr *QueryResult) GetChildRelOne(rel Relationship) *QueryResult {
	if !rel.Multi() {
		if v, ok := qr.ToOne[rel.Name()]; ok {
			return v
		}
	}
	panic("Can't get one on a multi relationship")
}

func (qr *QueryResult) GetChildRelMany(rel Relationship) []*QueryResult {
	if rel.Multi() {
		if v, ok := qr.ToMany[rel.Name()]; ok {
			return v
		}
	}
	panic("Can't get one on a multi relationship")
}

func (qr *QueryResult) SetChildRelMany(key string, qrs []*QueryResult) {
	if qr.ToMany == nil {
		qr.ToMany = map[string][]*QueryResult{key: qrs}
	} else {
		qr.ToMany[key] = qrs
	}
}

type ModelRef struct {
	InterfaceID ID
	AliasID     uuid.UUID
	I           Interface
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
	To ModelRef
	on RefRelationship
	jt joinType
}

func (j JoinOperation) String() string {
	return fmt.Sprintf("join: %v\t(%v)", j.To.I.Name(), j.To.AliasID)
}

type SetOpType int

const (
	or SetOpType = iota
	and
	not
)

type SetOperation struct {
	op       SetOpType
	Branches []Q
}

func (j JoinOperation) IsToOne() bool {
	return !j.on.rel.Multi()
}

func (j JoinOperation) Key() string {
	return j.on.rel.Name()
}

func (tx *holdTx) Ref(interfaceID ID) ModelRef {
	i, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		panic("Bad ref")
	}
	return ModelRef{interfaceID, uuid.New(), i}
}

type QueryClause func(*Q)

func (tx *holdTx) Query(model ModelRef, clauses ...QueryClause) Q {
	qb := initQB()
	qb.tx = tx
	qb.Root = &model
	for _, c := range clauses {
		c(&qb)
	}
	return qb
}

func LeftJoin(to ModelRef, on RefRelationship) QueryClause {
	return func(qb *Q) {
		outer := on.from
		j := JoinOperation{to, on, leftJoin}
		joinList, ok := qb.Joins[outer.AliasID]
		if ok {
			qb.Joins[outer.AliasID] = append(joinList, j)
		} else {
			qb.Joins[outer.AliasID] = []JoinOperation{j}
		}
	}
}

func Join(to ModelRef, on RefRelationship) QueryClause {
	return func(qb *Q) {
		outer := on.from
		j := JoinOperation{to, on, innerJoin}
		joinList, ok := qb.Joins[outer.AliasID]
		if ok {
			qb.Joins[outer.AliasID] = append(joinList, j)
		} else {
			qb.Joins[outer.AliasID] = []JoinOperation{j}
		}

	}
}

func Filter(ref ModelRef, m Matcher) QueryClause {
	return func(qb *Q) {
		matcherList, ok := qb.Filters[ref.AliasID]
		if ok {
			qb.Filters[ref.AliasID] = append(matcherList, m)
		} else {
			qb.Filters[ref.AliasID] = []Matcher{m}
		}
	}
}

func Aggregate(ref ModelRef, a Aggregation) QueryClause {
	return func(qb *Q) {
		qb.Aggregations[ref.AliasID] = a
	}
}

func Or(ref ModelRef, branches ...Q) QueryClause {
	return SetOpClause(ref, or, branches...)
}

func Intersection(ref ModelRef, branches ...Q) QueryClause {
	return SetOpClause(ref, and, branches...)
}

func Not(ref ModelRef, branches ...Q) QueryClause {
	return SetOpClause(ref, not, branches...)
}

func SetOpClause(ref ModelRef, op SetOpType, branches ...Q) QueryClause {
	return func(qb *Q) {
		sos, ok := qb.SetOps[ref.AliasID]
		so := SetOperation{op, branches}
		if ok {
			qb.SetOps[ref.AliasID] = append(sos, so)
		} else {
			qb.SetOps[ref.AliasID] = []SetOperation{so}
		}
	}
}

func (q Q) All() []*QueryResult {
	results := q.runBlockRoot(q.tx)
	return results
}

func (q Q) Records() []Record {
	results := q.runBlockRoot(q.tx)
	recs := []Record{}
	for _, r := range results {
		recs = append(recs, r.Record)
	}
	return recs
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

type Q struct {
	// TODO should we actually export these (vs exporting some setters)
	tx *holdTx

	// null if this isn't a Root QB
	Root *ModelRef

	// each of these are keyed by aliasID
	Aggregations map[uuid.UUID]Aggregation
	Joins        map[uuid.UUID][]JoinOperation
	Filters      map[uuid.UUID][]Matcher
	SetOps       map[uuid.UUID][]SetOperation
}

func (q Q) String() string {
	var ifID string
	if q.Root == nil {
		ifID = "subquery"
	} else {
		ifID = q.Root.InterfaceID.String()
	}
	return fmt.Sprintf(`Root: %v
Aggregations: %v
Joins: %v
Filters: %v
SetOps: %v`, ifID, q.Aggregations, q.Joins, q.Filters, q.SetOps)
}

func initQB() Q {
	return Q{
		Aggregations: map[uuid.UUID]Aggregation{},
		Filters:      map[uuid.UUID][]Matcher{},
		SetOps:       map[uuid.UUID][]SetOperation{},
		Joins:        map[uuid.UUID][]JoinOperation{}}
}

func Subquery(clauses ...QueryClause) Q {
	qb := initQB()
	for _, c := range clauses {
		c(&qb)
	}
	return qb
}
