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
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Aggregation int

const (
	Some Aggregation = iota
	None
	Every
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

type set map[string]void

type Selection struct {
	selecting bool
	fields    set
}

func (s Selection) String() string {
	if s.selecting {
		var keys []string
		for k := range s.fields {
			keys = append(keys, k)
		}
		return fmt.Sprintf("selection%v", keys)
	}
	return "selection[]"
}

type QueryResult struct {
	selects set
	Record  Record
	ToOne   map[string]*QueryResult
	ToMany  map[string][]*QueryResult
}

func (qr *QueryResult) HideAll() {
	if qr.selects == nil {
		qr.selects = make(set)
	}
}

func (qr *QueryResult) Show(field string) {
	_, err := qr.Record.Interface().AttributeByName(field)
	if err == nil {
		qr.selects[field] = void{}
	}
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

func (qr *QueryResult) Map() (map[string]interface{}, error) {
	if qr.Record == nil {
		return nil, nil
	}
	data := qr.Record.Map()
	if qr.selects != nil {
		for k, _ := range data {
			if _, ok := qr.selects[k]; !ok {
				delete(data, k)
			}
		}
	}
	for k, v := range qr.ToOne {
		data[k] = v
	}
	for k, v := range qr.ToMany {
		data[k] = v
	}
	return data, nil
}

func (qr *QueryResult) MarshalJSON() ([]byte, error) {
	if qr.Record == nil {
		return json.Marshal(nil)
	}
	data, err := qr.Map()
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func (qr *QueryResult) Get(name string) (val interface{}, err error) {
	if qr.Record == nil {
		return nil, nil
	}
	val, err = qr.Record.Get(name)
	if err == nil {
		return
	}
	val, ok := qr.ToOne[name]
	if ok {
		return val, nil
	}
	val, ok = qr.ToMany[name]
	if ok {
		return val, nil
	}
	return nil, errors.New("No such name")
}

func (qr *QueryResult) GetChildRelOne(rel Relationship) *QueryResult {
	if !rel.Multi() {
		if v, ok := qr.ToOne[rel.Name()]; ok {
			return v
		}
		return nil
	}
	panic("Can't get one on a multi relationship")
}

func (qr *QueryResult) GetChildRelMany(rel Relationship) []*QueryResult {
	if rel.Multi() {
		if v, ok := qr.ToMany[rel.Name()]; ok {
			return v
		}
		return nil
	}
	panic("Can't get one on a multi relationship")
}

func (qr *QueryResult) SetChildRelMany(key string, related []*QueryResult) {
	if qr.ToMany == nil {
		qr.ToMany = map[string][]*QueryResult{key: related}
	} else {
		qr.ToMany[key] = related
	}
}

func (qr *QueryResult) SetChildRelOne(key string, related *QueryResult) {
	if qr.ToOne == nil {
		qr.ToOne = map[string]*QueryResult{key: related}
	} else {
		qr.ToOne[key] = related
	}
}

type ModelRef struct {
	InterfaceID ID
	AliasID     uuid.UUID
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

func (jt joinType) String() string {
	switch jt {
	case leftJoin:
		return "leftJoin"
	case innerJoin:
		return "innerJoin"
	default:
		panic("Invalid joinType")
	}
}

type JoinOperation struct {
	To ModelRef
	on RefRelationship
	jt joinType
}

func (j JoinOperation) String() string {
	return fmt.Sprintf("join: %v on %v (%v)", j.To.InterfaceID, j.on.rel.Name(), j.To.AliasID)
}

type SetOpType int

const (
	or SetOpType = iota
	and
	not
)

type SetOperation struct {
	op SetOpType

	// do these really need to be whole new Qs?
	Branches []Q
}

type CaseOperation struct {
	Of ModelRef
}

func (c CaseOperation) String() string {
	return fmt.Sprintf("case: %v (%v)", c.Of.InterfaceID, c.Of.AliasID)
}

type Sort struct {
	AttributeName string
	Ascending     bool
}

func Limit(ref ModelRef, limit int) QueryClause {
	return func(qb *Q) {
		qb.Limits[ref.AliasID] = limit
	}
}

func Offset(offset int, ref ModelRef) QueryClause {
	return func(qb *Q) {
		qb.Offsets[ref.AliasID] = offset
	}
}

func Order(ref ModelRef, sorts []Sort) QueryClause {
	return func(qb *Q) {
		qb.Orderings[ref.AliasID] = sorts
	}

}

func Case(from, of ModelRef) QueryClause {
	return func(qb *Q) {
		co := CaseOperation{of}
		caseList, ok := qb.Cases[from.AliasID]
		if ok {
			qb.Cases[from.AliasID] = append(caseList, co)
		} else {
			qb.Cases[from.AliasID] = []CaseOperation{co}
		}
	}
}

func (j JoinOperation) IsToOne() bool {
	return !j.on.rel.Multi()
}

func (j JoinOperation) Key() string {
	return j.on.rel.Name()
}

func (tx *holdTx) Ref(interfaceID ID) ModelRef {
	return ModelRef{interfaceID, uuid.New()}
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

func Select(ref ModelRef, fields []string) QueryClause {
	//TODO verify it's a valid field name
	return func(qb *Q) {
		selection, ok := qb.Selections[ref.AliasID]
		if ok {
			for _, field := range fields {
				selection.fields[field] = void{}
			}
		} else {
			qb.Selections[ref.AliasID] = Selection{true, make(set)}
			for _, field := range fields {
				qb.Selections[ref.AliasID].fields[field] = void{}
			}

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
	rootNode := Plan(q)
	rIter, err := rootNode.ResultIter(q.tx, nil)
	if err != nil {
		panic(err)
	}
	newResults := []*QueryResult{}
	for rIter.Next() {
		qr := rIter.Value()
		newResults = append(newResults, qr)
	}
	if rIter.Err() != Done {
		panic(err)
	}

	return newResults
}

func (q Q) Records() []Record {
	results := q.All()
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
		err := fmt.Errorf("Expected one record but found many: %v", results)
		panic(err)
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
	Selections   map[uuid.UUID]Selection
	SetOps       map[uuid.UUID][]SetOperation
	Cases        map[uuid.UUID][]CaseOperation
	Orderings    map[uuid.UUID][]Sort
	Limits       map[uuid.UUID]int
	Offsets      map[uuid.UUID]int
}

func (q Q) String() string {
	var ifID, ifAlias string
	if q.Root == nil {
		ifID = "subquery"
		ifAlias = "subquery"
	} else {
		ifID = q.Root.InterfaceID.String()
		ifAlias = q.Root.AliasID.String()
	}
	return fmt.Sprintf(`Root: %v (%v)
Aggregations: %v
Joins: %v
Filters: %v
SetOps: %v
Cases: %v`, ifID, ifAlias, q.Aggregations, q.Joins, q.Filters, q.SetOps, q.Cases)
}

func initQB() Q {
	return Q{
		Aggregations: map[uuid.UUID]Aggregation{},
		Filters:      map[uuid.UUID][]Matcher{},
		Selections:   map[uuid.UUID]Selection{},
		SetOps:       map[uuid.UUID][]SetOperation{},
		Joins:        map[uuid.UUID][]JoinOperation{},
		Cases:        map[uuid.UUID][]CaseOperation{},
		Orderings:    map[uuid.UUID][]Sort{},
		Limits:       map[uuid.UUID]int{},
		Offsets:      map[uuid.UUID]int{},
	}
}

func (tx *holdTx) Subquery(clauses ...QueryClause) Q {
	qb := initQB()
	qb.tx = tx
	for _, c := range clauses {
		c(&qb)
	}
	return qb
}
