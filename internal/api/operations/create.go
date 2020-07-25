package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op CreateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	tx.Insert(op.Record)

	root := tx.Ref(op.Record.Interface().ID())
	parents := []*db.QueryResult{&db.QueryResult{Record: op.Record}}
	clauses := []db.QueryClause{}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, root, parents, clauses)
		if err != nil {
			return nil, err
		}
	}

	ids := []db.ID{op.Record.ID()}
	clauses = []db.QueryClause{db.Filter(root, db.IDIn(ids))}
	clauses = append(clauses, handleIncludes(tx, root, op.FindArgs.Include)...)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	if len(qrs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	tx.Commit()
	return qrs[0], nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	tx.Insert(op.Record)
	for _, parent := range parents {
		tx.Connect(parent.Record.ID(), op.Record.ID(), op.Relationship.ID())
	}
	for _, no := range op.Nested {
		err = no.ApplyNested(tx, root, parent, []*db.QueryResult{&db.QueryResult{Record: op.Record}}, clauses)
		if err != nil {
			return
		}
	}

	return nil
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, _ := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := getEdgeResults(parents, q.All())

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		rec := outs[0].Record
		for _, parent := range parents {
			tx.Connect(parent.Record.ID(), rec.ID(), op.Relationship.ID())
		}
	}
	return
}
