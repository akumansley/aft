package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

func (op CreateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	rec, err := buildRecordFromData(tx, op.ModelID, op.Data)
	if err != nil {
		return nil, err
	}
	tx.Insert(rec)

	root := tx.Ref(rec.InterfaceID())
	parents := []*db.QueryResult{&db.QueryResult{Record: rec}}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, parents)
		if err != nil {
			return nil, err
		}
	}

	ids := []db.ID{rec.ID()}
	clauses := []db.QueryClause{db.Filter(root, db.IDIn(ids))}
	clauses = append(clauses, handleIncludes(tx, root, op.FindArgs.Include)...)
	clauses = append(clauses, handleSelects(tx, root, op.FindArgs.Select)...)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	if len(qrs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	return qrs[0], nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	rec, err := buildRecordFromData(tx, op.Model.ID(), op.Data)
	if err != nil {
		return err
	}
	tx.Insert(rec)
	for _, parent := range parents {
		err = tx.Connect(parent.Record.ID(), rec.ID(), op.Relationship.ID())
		if err != nil {
			return err
		}
	}
	for _, no := range op.Nested {
		err = no.ApplyNested(tx, parent, []*db.QueryResult{&db.QueryResult{Record: rec}})
		if err != nil {
			return
		}
	}

	return nil
}

func buildRecordFromData(tx db.RWTx, modelID db.ID, data map[string]interface{}) (db.Record, error) {
	m, err := tx.Schema().GetInterfaceByID(modelID)
	if err != nil {
		return nil, err
	}
	u := uuid.New()
	rec, err := tx.MakeRecord(modelID)
	if err != nil {
		return nil, err
	}
	err = rec.Set("id", u)
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		a, err := m.AttributeByName(tx, k)
		if err != nil {
			return nil, err
		}
		err = a.Set(tx, rec, v)
		if err != nil {
			return nil, err
		}
	}
	return rec, nil
}
