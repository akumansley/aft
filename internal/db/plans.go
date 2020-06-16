package db

import (
	"github.com/google/uuid"
)

type framemaker struct {
	capacity int
	amap     map[uuid.UUID]int
}

func (fm framemaker) ix(aliasID uuid.UUID) int {
	return fm.amap[aliasID]
}

func (fm framemaker) frame(rec Record, ix int) frame {
	sl := make([]Record, fm.capacity)
	return frame{entries: sl}
}

type intermediate interface {
	ID() uuid.UUID
}

type frame struct {
	entries []intermediate
}

func (f frame) merge(other frame) {
	for ix := range other.entries {
		v, ok := other.entries[ix]
		if ok {
			f.entries[ix] = v
		}
	}
}

func newFramemaker(q *Query) *framemaker {
	numScans := len(q.aset)
	fm := framemaker{capacity: numScans}
	ix := 0
	for k := range q.aset {
		fm.amap[k] = ix
		ix++
	}
	return &fm
}

type plannode interface {
	iter(Tx) *frameiter
	lookup(Tx, string, uuid.UUID) frame
}

type lookup struct {
	modelID uuid.UUID
	matcher Matcher
	frameIx int
	fm      *framemaker
}

func (l lookup) iter(tx Tx) iterator {
	panic("shouldn't call iter on a lookup")
}

func (l lookup) lookup(tx Tx, key string, id uuid.UUID) frame {
	val, _ := tx.FindOne(l.modelID, And(Eq(key, id), l.matcher))
	return l.fm.frame(val, l.frameIx)
}

type scan struct {
	modelID uuid.UUID
	matcher Matcher
	frameIx int
	fm      *framemaker
}

func (s scan) iter(tx Tx) iterator {
	recs := tx.FindMany(s.modelID, s.matcher)
	var frames []frame
	for _, r := range recs {
		frames = append(frames, s.fm.frame(r, s.frameIx))
	}
	return &frameiter{frames: frames}
}

func (s scan) lookup(tx Tx, key string, aliasID, id uuid.UUID) frame {
	panic("shouldn't call lookup on a scan")
}

// a to-one loop join
type loopjoin struct {
	left, right plannode
	on          OnPredicate
	fm          *framemaker
	results     []frame
}

func (l loopjoin) lookup(tx Tx, key string, aliasID, id uuid.UUID) frame {
	panic("shouldn't call lookup on a loopjoin")
}

func (l loopjoin) iter(tx Tx) iterator {
	lFix := l.fm.ix(l.on.left.aliasID)
	rFix := l.fm.ix(l.on.right.aliasID)
	var results []frame

	if l.on.b.RelType() == BelongsTo {
		fi := l.left.iter(tx)
		for fi.Next() {
			lv := fi.Value()
			fk := v.GetFK(b.Name())
			rv := l.right.lookup(tx, "id", fk)
			lv.merge(rv)
			results = append(results, lv)
		}
	} else if l.on.b.RelType() == HasOne {
		fi := l.left.iter(tx)
		key := b.Dual().Name()
		for fi.Next() {
			lv := fi.Value()
			rv := l.right.lookup(tx, key, lv.ID())
			lv.merge(rv)
			results = append(results, lv)
		}
	}
	l.results = results
	return &frameiter{frames: results}
}

type hashjoinmany struct {
	left, right plannode
	on          OnPredicate
	agg         Aggregation
}

func (l hashjoinmany) lookup(tx Tx, key string, aliasID, id uuid.UUID) frame {
	panic("shouldn't call lookup on a hashjoin")
}

func (j hashjoinmany) iter(tx Tx) iterator {
	b := j.on.b
	fk := b.Dual().Name()
	lFix := j.fm.ix(l.on.left.aliasID)
	rFix := j.fm.ix(l.on.right.aliasID)
	var results []frame

	// join
	if j.on.b.RelType() == BelongsTo {
		fi := j.left.iter(tx)
		for fi.Next() {
			lv := fi.Value()
		}
		// build hash ^ to list
		fk := v.GetFK(b.Name())
		rv := j.right.lookup(tx, "id", fk)
		lv.merge(rv)
		results = append(results, lv)
	}

	// aggregate
	it, err := j.iter(tx)
	if err != nil {
		return nil, err
	}
	hash := map[uuid.UUID][]Record{}
	for it.Next() {
		v := it.Value()
		k := v.GetFK(fk)
		ls, ok := hash[k]
		if ok {
			hash[k] = append(ls, v)
		} else {
			hash[k] = []Record{v}
		}
	}
	if it.Err() != nil {
		return nil, it.Err()
	}
	return &groupiter{groups: hash}, nil

}

type hashsetop struct {
	left, right       plannode
	leftRef, rightRef ModelRef
	op                setoperation
}

func (l hashsetop) lookup(tx Tx, key string, aliasID, id uuid.UUID) frame {
	panic("not implemented")
}

func (l hashsetop) iter(tx Tx) frameiter {
	panic("not implemented")
}

func plan(q *Query) plannode {
	fm := newFramemaker(q)
	return planRelation(fm, q, root{}, q.r)
}

func planRelation(fm *framemaker, q *Query, parent relation, r relation) plannode {
	switch r.(type) {
	case table:
		t := r.(table)
		if parent.shouldScan(t) {
			return scan{modelID: t.ref.modelID, fm: fm}
		} else {
			return lookup{modelID: t.ref.modelID, fm: fm}
		}
	case joinone:
		j1 := r.(joinone)
		left := planRelation(fm, q, j1, j1.left)
		right := planRelation(fm, q, j1, j1.right)
		return loopjoin{left, right, j1.on}
	case joinmany:
		jm := r.(joinmany)
		left := planRelation(fm, q, jm, jm.left)
		right := planRelation(fm, q, jm, jm.right)
		return hashjoinmany{left, right, jm.on, jm.agg}
	case setop:
		so := r.(setop)
		left := planRelation(fm, q, so, so.left)
		right := planRelation(fm, q, so, so.right)
		return hashsetop{left, right, so.leftRef, so.rightRef, so.op}
	}
	return nil
}
