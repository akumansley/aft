package operations

import (
	"awans.org/aft/internal/db"
)

type Inclusion struct {
	Relationship   db.Relationship
	NestedFindMany FindManyArgs
}

type Include struct {
	Includes []Inclusion
}

func (i Include) One(tx db.Tx, m db.ID, rec db.Record) (*db.QueryResult, error) {
	recs := []db.Record{rec}
	q := buildIncQuery(tx, m, recs, i)
	qrs := q.All()
	if len(qrs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	return qrs[0], nil
}

func buildIncQuery(tx db.Tx, m db.ID, recs []db.Record, i Include) db.Q {
	ids := []db.ID{}
	for _, r := range recs {
		ids = append(ids, r.ID())
	}

	root := tx.Ref(m)
	q := tx.Query(root)
	q = q.Filter(root, db.IDIn(ids))
	qb := q.AsBlock()
	qb = handleIncludes(tx, qb, root, i)
	q.SetMainBlock(qb)
	return q
}

func handleIncludes(tx db.Tx, qb db.QBlock, parent db.ModelRef, i Include) db.QBlock {
	for _, inclusion := range i.Includes {
		qb = handleInclusion(tx, qb, parent, inclusion)
	}
	return qb
}

func handleInclusion(tx db.Tx, q db.QBlock, parent db.ModelRef, i Inclusion) db.QBlock {
	child := tx.Ref(i.Relationship.Target().ID())
	qb := q.LeftJoin(child, parent.Rel(i.Relationship))
	if i.Relationship.Multi() {
		qb.Aggregate(child, db.Include)
	}
	clauses = append(clauses, handleFindMany(tx, child, i.NestedFindMany)...)
	return clauses
}
