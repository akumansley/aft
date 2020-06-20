package db

import (
	"github.com/google/uuid"
)

// Utility methods
// filterEmtpy -- copies a []*QR into a new slice, removing any that are isEmpty()

func filterEmpty(results []*QueryResult) []*QueryResult {
	filtered := []*QueryResult{}
	for _, qr := range results {
		if !qr.isEmpty() {
			filtered = append(filtered, qr)
		}
	}
	return filtered
}

func copyResultsShallow(results []*QueryResult) []*QueryResult {
	var copied []*QueryResult
	for _, r := range results {
		copied = append(copied, &QueryResult{Record: r.Record})
	}
	return copied
}

func applyMatcher(results []*QueryResult, matcher Matcher) []*QueryResult {
	for _, result := range results {
		if !result.isEmpty() {
			match, _ := matcher.Match(result.Record)
			if !match {
				result.Empty()
			}
		}
	}
	return results
}

// Entrypoint

func (qb QBlock) runBlockRoot(tx Tx) []*QueryResult {
	matchers := qb.sargs[qb.root.aliasID]
	outer := qb.performScan(tx, qb.root.modelID, And(matchers...))
	results := qb.runBlock(tx, outer, qb.root.aliasID)
	results = filterEmpty(results)
	return results
}

func (qb QBlock) runBlockNested(tx Tx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	matchers, ok := qb.sargs[aliasID]
	if ok {
		outer = applyMatcher(outer, And(matchers...))
	}
	return qb.runBlock(tx, outer, aliasID)
}

func (qb QBlock) runBlock(tx Tx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	results := qb.performJoins(tx, outer, aliasID)
	results = qb.performSetOps(tx, results, aliasID)
	return results
}

func (qb QBlock) performSetOps(tx Tx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	for _, s := range qb.setops[aliasID] {
		outer = qb.performSetOp(tx, outer, s, aliasID)
	}
	return outer
}

func orResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {

	for i, o := range original {
		any := false
		for j := range set {
			r := set[j][i]
			if !r.isEmpty() {
				any = true
				break
			}
		}
		if !any {
			o.Empty()
		}
	}
	return original
}

func andResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {
	for i, o := range original {
		all := true
		for j := range set {
			r := set[j][i]
			if r.isEmpty() {
				all = false
			}
		}
		if !all {
			o.Empty()
		}
	}
	return original
}

func notResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {
	for i, o := range original {
		any := false
		for j := range set {
			r := set[j][i]
			if !r.isEmpty() {
				any = false
			}
		}
		if any {
			o.Empty()
		}
	}
	return original
}

func (qb QBlock) performSetOp(tx Tx, outer []*QueryResult, op setop, aliasID uuid.UUID) []*QueryResult {
	original := copyResultsShallow(outer)
	var set [][]*QueryResult
	for _, b := range op.branches {
		branchCopy := copyResultsShallow(outer)
		branchResults := b.runBlockNested(tx, branchCopy, aliasID)
		set = append(set, branchResults)
	}
	switch op.op {
	case or:
		return orResults(original, set)
	case and:
		return andResults(original, set)
	case not:
		return notResults(original, set)
	default:
		panic("invalid set op")
	}
}

func (qb QBlock) performScan(tx Tx, modeID uuid.UUID, matcher Matcher) []*QueryResult {
	recs, _ := tx.FindMany(qb.root.modelID, matcher)
	var results []*QueryResult
	for _, rec := range recs {
		results = append(results, &QueryResult{Record: rec})
	}
	return results
}

func (qb QBlock) performJoins(tx Tx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	for _, j := range qb.joins[aliasID] {
		toOne := j.IsToOne()

		if toOne {
			outer = qb.performJoinOne(tx, outer, j)
		} else {
			outer = qb.performJoinMany(tx, outer, j)
		}
	}
	return outer
}

func (qb QBlock) performJoinOne(tx Tx, outer []*QueryResult, j join) []*QueryResult {
	var inner []*QueryResult
	key := j.Key()
	matchers := qb.sargs[j.to.aliasID]

	for _, r := range outer {
		qr := getRelatedOne(tx, r.Record, j, And(matchers...))
		inner = append(inner, qr)
	}

	inner = qb.performJoins(tx, inner, j.to.aliasID)

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
		thisFK, _ := rec.GetFK(b.Name())
		hit, _ := tx.FindOne(d.ModelID(), And(Eq("id", thisFK), matcher))
		return &QueryResult{Record: hit}
	}
	panic("invalid join")
}

func (qb QBlock) performJoinManySomeOrInclude(tx Tx, outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	key := j.Key()
	matchers := qb.sargs[j.to.aliasID]
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, And(matchers...))
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
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
	qb.performJoins(tx, uniqValues, j.to.aliasID)

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

func (qb QBlock) performJoinManyNone(tx Tx, outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, nil)
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
		}
	}

	matchers := qb.sargs[j.to.aliasID]
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
	qb.performJoins(tx, uniqValues, j.to.aliasID)

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

func (qb QBlock) performJoinManyEvery(tx Tx, outer []*QueryResult, j join, a Aggregation) []*QueryResult {
	key := j.Key()
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, nil)
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
		}
	}

	matchers := qb.sargs[j.to.aliasID]
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
	qb.performJoins(tx, uniqValues, j.to.aliasID)

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
func (qb QBlock) performJoinMany(tx Tx, outer []*QueryResult, j join) []*QueryResult {
	agg, ok := qb.aggregations[j.to.aliasID]
	if !ok {
		agg = Include
	}
	switch agg {
	case Some, Include:
		return qb.performJoinManySomeOrInclude(tx, outer, j, agg)
	case Every:
		return qb.performJoinManyEvery(tx, outer, j, agg)
	case None:
		return qb.performJoinManyNone(tx, outer, j, agg)
	}
	panic("not implemented")
}

func getRelatedMany(tx Tx, rec Record, j join, matcher Matcher) []*QueryResult {
	if rec == nil {
		panic("can't get related many of nil")
	}
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
		hits, _ := tx.FindMany(d.ModelID(), matcher)
		var results []*QueryResult
		for _, h := range hits {
			results = append(results, &QueryResult{Record: h})
		}
		return results
	}
	panic("invalid join")
}
