package operations

import (
	"awans.org/aft/internal/db"
)

type Inclusion struct {
	Relationship db.Relationship
	Where        Where
}

type Include struct {
	Includes []Inclusion
}

func (i Include) Resolve(tx db.Tx, m db.ID, recs []db.Record) []*db.QueryResult {
	q := buildIncQuery(tx, m, recs, i)
	return q.All()
}

func (i Include) ResolveOne(tx db.Tx, m db.ID, rec db.Record) *db.QueryResult {
	recs := []db.Record{rec}
	qrs := i.Resolve(tx, m, recs)
	if len(qrs) != 1 {
		panic("Resolve single include returned non-1 results")
	}
	return qrs[0]
}

func buildIncQuery(tx db.Tx, m db.ID, recs []db.Record, i Include) db.Q {
	ids := []db.ID{}
	for _, r := range recs {
		ids = append(ids, r.ID())
	}

	root := tx.Ref(m)
	f := db.Filter(root, db.IDIn(ids))
	clauses := []db.QueryClause{f}
	for _, inclusion := range i.Includes {
		clauses = append(clauses, handleInclusion(tx, root, inclusion)...)
	}

	q := tx.Query(root, clauses...)
	return q
}

func handleInclusion(tx db.Tx, parent db.ModelRef, i Inclusion) (clauses []db.QueryClause) {
	child := tx.Ref(i.Relationship.Target().ID())
	j := db.LeftJoin(child, parent.Rel(i.Relationship))
	clauses = append(clauses, j)
	if i.Relationship.Multi() {
		a := db.Aggregate(child, db.Include)
		clauses = append(clauses, a)
	}
	clauses = append(clauses, handleWhere(tx, child, i.Where)...)
	return clauses
}
