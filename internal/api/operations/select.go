package operations

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
)

type Select struct {
	Selecting bool
	Fields    api.Set
	Selects   []Selection
}

type Selection struct {
	Relationship   db.Relationship
	NestedFindMany FindArgs
}

func handleSelects(tx db.Tx, parent db.ModelRef, s Select) []db.QueryClause {
	clauses := []db.QueryClause{}
	if s.Selecting {
		var fields []string
		for k, _ := range s.Fields {
			fields = append(fields, k)
		}
		clauses = append(clauses, db.Select(parent, fields))
	}
	for _, selection := range s.Selects {
		clauses = append(clauses, handleSelection(tx, parent, selection)...)
	}
	return clauses
}

func handleSelection(tx db.Tx, parent db.ModelRef, s Selection) (clauses []db.QueryClause) {
	child := tx.Ref(s.Relationship.Target(tx).ID())
	j := db.LeftJoin(child, parent.Rel(s.Relationship))
	clauses = append(clauses, j)
	if s.Relationship.Multi() {
		a := db.Aggregate(child, db.Include)
		clauses = append(clauses, a)
	}
	clauses = append(clauses, handleFindMany(tx, child, s.NestedFindMany)...)
	return clauses
}
