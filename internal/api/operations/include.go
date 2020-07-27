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
	f := db.Filter(root, db.IDIn(ids))
	clauses := []db.QueryClause{f}
	handleIncludes(tx, root, i)
	q := tx.Query(root, clauses...)
	return q
}

func handleIncludes(tx db.Tx, parent db.ModelRef, i Include) (clauses []db.QueryClause) {
	for _, inclusion := range i.Includes {
		clauses = append(clauses, handleInclusion(tx, parent, inclusion)...)
	}
	return
}

func handleInclusion(tx db.Tx, parent db.ModelRef, i Inclusion) (clauses []db.QueryClause) {
	child := tx.Ref(i.Relationship.Target().ID())
	j := db.LeftJoin(child, parent.Rel(i.Relationship))
	clauses = append(clauses, j)
	if i.Relationship.Multi() {
		a := db.Aggregate(child, db.Include)
		clauses = append(clauses, a)
	}
	clauses = append(clauses, handleFindMany(tx, child, i.NestedFindMany)...)
	return clauses
}
