package api

import (
	"awans.org/aft/internal/db"
	"encoding/json"
)

type Inclusion struct {
	Relationship db.Relationship
	Where        Where
}

type Include struct {
	Includes []Inclusion
}

type IncludeResult struct {
	Record         db.Record
	SingleIncludes map[string]db.Record
	MultiIncludes  map[string][]db.Record
}

func (ir IncludeResult) MarshalJSON() ([]byte, error) {
	data := ir.Record.Map()
	for k, v := range ir.SingleIncludes {
		data[k] = v
	}
	for k, v := range ir.MultiIncludes {
		data[k] = v
	}
	return json.Marshal(data)
}

func (i Include) Resolve(tx db.Tx, m db.ModelID, recs []db.Record) []*db.QueryResult {
	q := buildIncQuery(tx, m, recs, i)
	return q.All()
}

func (i Include) ResolveOne(tx db.Tx, m db.ModelID, rec db.Record) *db.QueryResult {
	recs := []db.Record{rec}
	qrs := i.Resolve(tx, m, recs)
	if len(qrs) != 1 {
		panic("Resolve single include returned non-1 results")
	}
	return qrs[0]
}

func buildIncQuery(tx db.Tx, m db.ModelID, recs []db.Record, i Include) db.Q {
	ids := []db.ID{}
	for _, r := range recs {
		ids = append(ids, r.ID())
	}

	root := tx.Ref(m)
	q := tx.Query(root)
	q = q.Filter(root, db.IDIn(ids))
	qb := q.AsBlock()
	for _, inclusion := range i.Includes {
		qb = handleInclusion(tx, root, qb, inclusion)
	}
	q.SetMainBlock(qb)
	return q
}

func handleInclusion(tx db.Tx, parent db.ModelRef, q db.QBlock, i Inclusion) db.QBlock {
	child := tx.Ref(i.Relationship.Target.ID)
	qb := q.LeftJoin(child, parent.Rel(i.Relationship))
	if i.Relationship.Multi {
		qb.Aggregate(child, db.Include)
	}
	qb = handleWhere(tx, qb, child, i.Where)
	return qb
}
