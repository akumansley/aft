package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

type Case struct {
	Entries []CaseEntry
}

type CaseEntry struct {
	ModelID db.ID
	Include Include
	Select  Select
}

func (c Case) String() string {
	return fmt.Sprintf("case{%v}\n", c.Entries)
}

func (e CaseEntry) String() string {
	return fmt.Sprintf("entry{%v %v %v}\n", e.ModelID, e.Include, e.Select)
}

func handleCase(tx db.Tx, parent db.ModelRef, c Case) []db.QueryClause {
	clauses := []db.QueryClause{}
	for _, entry := range c.Entries {
		clauses = append(clauses, handleCaseEntry(tx, parent, entry)...)
	}
	return clauses
}

func handleCaseEntry(tx db.Tx, parent db.ModelRef, e CaseEntry) (clauses []db.QueryClause) {
	// it's a new aliasID, even though it's the same physical "scan"
	cased := tx.Ref(e.ModelID)
	clause := db.Case(parent, cased)

	clauses = append(clauses, clause)

	clauses = append(clauses, handleIncludes(tx, cased, e.Include)...)
	clauses = append(clauses, handleSelects(tx, cased, e.Select)...)

	return clauses
}
