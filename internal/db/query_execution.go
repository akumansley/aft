package db

import (
	"encoding/json"
	"github.com/google/uuid"
)

func filterEmpty(results []*QueryResult) []*QueryResult {
	filtered := []*QueryResult{}
	for _, qr := range results {
		if !qr.isEmpty() {
			filtered = append(filtered, qr)
		}
	}
	return filtered
}
func (q Q) performScan(modeID uuid.UUID, matcher Matcher) []*QueryResult {
	recs := q.tx.FindMany(q.root.modelID, matcher)
	var results []*QueryResult
	for _, rec := range recs {
		results = append(results, &QueryResult{Record: rec})
	}
	return results
}

func (q Q) performJoins(outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	for _, j := range q.joins[aliasID] {
		toOne := j.IsToOne()

		if toOne {
			outer = q.performJoinOne(outer, j)
		} else {
			outer = q.performJoinMany(outer, j)
		}
	}
	return outer
}

func (q Q) performJoinOne(outer []*QueryResult, j join) []*QueryResult {
	var inner []*QueryResult
	key := j.Key()
	matchers := q.sargs[j.to.aliasID]

	for _, r := range outer {
		qr := getRelatedOne(q.tx, r.Record, j, And(matchers...))
		inner = append(inner, qr)
	}

	inner = q.performJoins(inner, j.to.aliasID)

	if j.jt == innerJoin {
		// inner join
		for i := range outer {
			if !inner[i].isEmpty() {
				outer[i].ToOne[key] = inner[i]
			} else {
				outer[i].Empty()
			}
		}
		return outer
	} else {
		// left join
		for i := range outer {
			if !inner[i].isEmpty() {
				outer[i].ToOne[key] = inner[i]
			}
		}
		return outer
	}
}

func getRelatedOne(tx Tx, rec Record, j join, matcher Matcher) *QueryResult {
	b := j.on.b
	d := b.Dual()
	id := rec.ID()
	switch b.RelType() {
	case HasOne:
		// FK on the other side
		hit, _ := tx.FindOne(d.ModelID(), And(EqFK(d.Name(), id), matcher))
		return &QueryResult{Record: hit}
	case BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(b.Name())
		hit, _ := tx.FindOne(d.ModelID(), And(Eq("id", thisFK), matcher))
		return &QueryResult{Record: hit}
	}
	panic("invalid join")
}

func (q Q) performJoinManySomeOrInclude(outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	key := j.Key()
	matchers := q.sargs[j.to.aliasID]
	var inner [][]*QueryResult
	for _, r := range outer {
		qr := getRelatedMany(q.tx, r.Record, j, And(matchers...))
		inner = append(inner, qr)
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[uuid.UUID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins
	// we're passing pointers so stuf'll get modified in-place in the dict
	q.performJoins(uniqValues, j.to.aliasID)

	// merge 'em back, filtering if none made it back
	for i := range outer {
		joinedSet := inner[i]
		var populatedJoinedSet []*QueryResult
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				populatedJoinedSet = append(populatedJoinedSet, populated)
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		isEmpty := true
		for _, v := range populatedJoinedSet {
			if !v.isEmpty() {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			// if this is a Some aggregation
			// blank out the parent record
			if a == Some {
				outer[i].Empty()
			}
		} else {
			dict := outer[i].ToMany
			if dict != nil {
				dict[key] = populatedJoinedSet
			} else {
				outer[i].ToMany = map[string][]*QueryResult{key: populatedJoinedSet}
			}
		}

	}
	return outer
}

func (q Q) performJoinManyNone(outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	var inner [][]*QueryResult
	for _, r := range outer {
		qr := getRelatedMany(q.tx, r.Record, j, nil)
		inner = append(inner, qr)
	}

	matchers := q.sargs[j.to.aliasID]
	matcher := And(matchers...)

	// apply the local filtering criteria
	// eagerly so we can maybe avoid doing some extra joins
	var filtered [][]*QueryResult
	for _, group := range inner {
		none := true
		for _, result := range group {
			match, _ := matcher.Match(result.Record)
			if match {
				none = false
				break
			}
		}
		if none {
			filtered = append(filtered, group)
		} else {
			filtered = append(filtered, []*QueryResult{})
		}
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[uuid.UUID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins
	// we're passing pointers so stuff'll get modified in-place in the dict
	q.performJoins(uniqValues, j.to.aliasID)

	// merge 'em back, filtering if any make it back
	for i := range outer {
		joinedSet := inner[i]
		none := true
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				if !populated.isEmpty() {
					none = false
					break
				}
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		if !none {
			// if this is an None aggregation
			// blank out the parent record
			if a == None {
				outer[i].Empty()
			}
		}
	}
	return outer
}

func (q Q) performJoinManyEvery(outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	key := j.Key()
	var inner [][]*QueryResult
	for _, r := range outer {
		qr := getRelatedMany(q.tx, r.Record, j, nil)
		inner = append(inner, qr)
	}

	matchers := q.sargs[j.to.aliasID]
	matcher := And(matchers...)

	// apply the local filtering criteria
	var filtered [][]*QueryResult
	// eagerly so we can maybe avoid doing some extra joins
	for _, group := range inner {
		every := true
		for _, result := range group {
			match, _ := matcher.Match(result.Record)
			if !match {
				every = false
				break
			}
		}
		if every {
			filtered = append(filtered, group)
		} else {
			filtered = append(filtered, []*QueryResult{})
		}
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[uuid.UUID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins
	// we're passing pointers so stuf'll get modified in-place in the dict
	q.performJoins(uniqValues, j.to.aliasID)

	// merge 'em back, filtering if any didn't make it back
	for i := range outer {
		joinedSet := inner[i]
		every := true
		var populatedJoinedSet []*QueryResult
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				if populated.isEmpty() {
					every = false
					break
				}
				populatedJoinedSet = append(populatedJoinedSet, populated)
			} else {
				every = false
				break
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		if !every {
			// if this is an Every aggregation
			// blank out the parent record
			if a == Every {
				outer[i].Empty()
			}
		} else {
			dict := outer[i].ToMany
			if dict != nil {
				dict[key] = populatedJoinedSet
			} else {
				outer[i].ToMany = map[string][]*QueryResult{key: populatedJoinedSet}
			}
		}

	}
	return outer
}

// returns QueryResults for just the right half of this one join
func (q Q) performJoinMany(outer []*QueryResult, j join) []*QueryResult {
	agg, ok := q.aggregations[j.to.aliasID]
	if !ok {
		agg = Include
	}
	switch agg {
	case Some, Include:
		return q.performJoinManySomeOrInclude(outer, j, agg)
	case Every:
		return q.performJoinManyEvery(outer, j, agg)
	case None:
		return q.performJoinManyNone(outer, j, agg)
	}
	panic("not implemented")
}

func getRelatedMany(tx Tx, rec Record, j join, matcher Matcher) []*QueryResult {
	b := j.on.b
	d := b.Dual()
	id := rec.ID()
	if matcher != nil {
		matcher = And(EqFK(d.Name(), id), matcher)
	} else {
		matcher = EqFK(d.Name(), id)
	}
	switch b.RelType() {
	case HasMany:
		// FK on the other side
		hits := tx.FindMany(d.ModelID(), matcher)
		var results []*QueryResult
		for _, h := range hits {
			results = append(results, &QueryResult{Record: h})
		}
		return results
	}
	panic("invalid join")
}
