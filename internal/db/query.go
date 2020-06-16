// user = db.Ref(modelID)
// db.Query(user).Join(post, user.Rel["posts"]).Filter(post, db.Eq("foo", "bar")).All()

package db

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type QueryResult struct {
	Record Record
	ToOne  map[string]QueryResult
	ToMany map[string][]QueryResult
}

func (qr QueryResult) Empty() bool {
	return qr.Record == nil
}

func (qr QueryResult) MarshalJSON() ([]byte, error) {
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
	left joinType = iota
	inner
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
	root  ModelRef
	joins map[uuid.UUID][]join
	tx    *holdTx
	sargs map[uuid.UUID][]Matcher
}

func (tx *holdTx) Ref(modelID uuid.UUID) ModelRef {
	model, _ := tx.GetModelByID(modelID)
	return ModelRef{modelID, uuid.New(), model}
}

func (tx *holdTx) Query(model ModelRef) Q {
	return Q{root: model, tx: tx,
		sargs: map[uuid.UUID][]Matcher{},
		joins: map[uuid.UUID][]join{}}
}

func (q Q) Join(to ModelRef, on RefBinding) Q {
	outer := on.from
	j := join{to, on, inner}
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

func (q Q) All() []QueryResult {
	var results []QueryResult
	matchers := q.sargs[q.root.aliasID]
	results = q.performScan(q.root.modelID, And(matchers...))

	results = q.performJoins(results, q.root.aliasID)
	return results
}

func (q Q) performScan(modeID uuid.UUID, matcher Matcher) []QueryResult {
	recs := q.tx.FindMany(q.root.modelID, matcher)
	var results []QueryResult
	for _, rec := range recs {
		results = append(results, QueryResult{Record: rec})
	}
	return results
}

func (q Q) performJoins(outer []QueryResult, aliasID uuid.UUID) []QueryResult {
	for _, j := range q.joins[aliasID] {
		toOne := j.IsToOne()
		key := j.Key()

		if toOne {
			inner := q.performJoinOne(outer, j)
			inner = q.performJoins(inner, j.to.aliasID)
			for i := range outer {
				if inner[i].Empty() {
					outer = append(outer[i:], outer[:i+1]...)
				} else {
					outer[i].ToOne[key] = inner[i]
				}
			}
		} else {
			inner := q.performJoinMany(outer, j)

			// first filter out anything that didn't survive the join
			fmt.Printf("toFilter: %v\n", inner)
			var filteredouter []QueryResult
			var filteredinner [][]QueryResult
			for i := range outer {
				qrList := inner[i]
				fmt.Printf("ok: %v\n", qrList)
				if len(qrList) > 0 {
					filteredouter = append(filteredouter, outer[i])
					filteredinner = append(filteredinner, inner[i])
				}
			}
			outer = filteredouter
			inner = filteredinner

			// to prevent explosion, we first merge by unique records
			// and then expand out
			uniq := map[uuid.UUID]QueryResult{}
			for _, group := range inner {
				for _, result := range group {
					uniq[result.Record.ID()] = result
				}
			}

			var uniqValues []QueryResult
			for _, uniqVal := range uniq {
				uniqValues = append(uniqValues, uniqVal)
			}
			uniqValues = q.performJoins(uniqValues, j.to.aliasID)
			for _, uniqVal := range uniqValues {
				uniq[uniqVal.Record.ID()] = uniqVal
			}

			for i := range outer {
				joinedSet := inner[i]
				var populatedJoinedSet []QueryResult
				for _, joined := range joinedSet {
					populated := uniq[joined.Record.ID()]
					populatedJoinedSet = append(populatedJoinedSet, populated)
				}
				dict := outer[i].ToMany
				if dict != nil {
					dict[key] = populatedJoinedSet
				} else {
					outer[i].ToMany = map[string][]QueryResult{key: populatedJoinedSet}
				}
			}
		}
	}
	return outer
}

// returns QueryResults for just the right half of this one join
// one to one with the input
func (q Q) performJoinOne(results []QueryResult, j join) []QueryResult {
	var outer []QueryResult
	matchers := q.sargs[j.to.aliasID]

	for _, r := range results {
		qr := getRelatedOne(q.tx, r.Record, j, And(matchers...))
		outer = append(outer, qr)
	}
	return outer
}

func getRelatedOne(tx Tx, rec Record, j join, matcher Matcher) QueryResult {
	b := j.on.b
	d := b.Dual()
	id := rec.ID()
	switch b.RelType() {
	case HasOne:
		// FK on the other side
		hit, _ := tx.FindOne(d.ModelID(), And(EqFK(d.Name(), id), matcher))
		return QueryResult{Record: hit}
	case BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(b.Name())
		hit, _ := tx.FindOne(d.ModelID(), And(Eq("id", thisFK), matcher))
		return QueryResult{Record: hit}
	}
	panic("invalid join")
}

// returns QueryResults for just the right half of this one join
func (q Q) performJoinMany(results []QueryResult, j join) [][]QueryResult {
	matchers := q.sargs[j.to.aliasID]
	var outer [][]QueryResult
	for _, r := range results {
		qr := getRelatedMany(q.tx, r.Record, j, And(matchers...))
		outer = append(outer, qr)
	}
	return outer
}

func getRelatedMany(tx Tx, rec Record, j join, matcher Matcher) []QueryResult {
	b := j.on.b
	d := b.Dual()
	id := rec.ID()
	switch b.RelType() {
	case HasMany:
		// FK on the other side
		hits := tx.FindMany(d.ModelID(), And(EqFK(d.Name(), id), matcher))
		var results []QueryResult
		for _, h := range hits {
			results = append(results, QueryResult{Record: h})
		}
		return results
	}
	panic("invalid join")
}
