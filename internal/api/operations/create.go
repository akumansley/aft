package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op CreateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	rec, err := buildRecordFromData(tx, op.ModelID, op.Data)
	if err != nil {
		return nil, err
	}
	tx.Insert(rec)

	root := tx.Ref(rec.Interface().ID())
	parents := []*db.QueryResult{&db.QueryResult{Record: rec}}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return nil, err
		}
	}

	ids := []db.ID{rec.ID()}
	clauses := []db.QueryClause{db.Filter(root, db.IDIn(ids))}
	clauses = append(clauses, handleIncludes(tx, root, op.FindArgs.Include)...)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	if len(qrs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	tx.Commit()
	return qrs[0], nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx) (err error) {
	tx.Insert(op.Record)
	tx.Connect(op.Relationship.Source().ID(), op.Record.ID(), op.Relationship.ID())

	for _, no := range op.Nested {
		err = no.ApplyNested(tx)
		if err != nil {
			return
		}
	}

	return nil
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	out := q.All()

	if len(out) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(out) == 1 {
		rec := out[0].Record
		tx.Connect(op.Relationship.Source().ID(), rec.ID(), op.Relationship.ID())
	}
	return
}

func buildRecordFromData(tx db.RWTx, modelID db.ID, data map[string]interface{}) (db.Record, error) {
	m, err := tx.Schema().GetInterfaceByID(modelID)
	if err != nil {
		return nil, err
	}
	rec := db.NewRecord(m)
	for k, v := range data {
		rec.Set(k, v)
	}
	return rec, nil
}
