package operations

import (
	"awans.org/aft/internal/db"
)

type FieldCriterion struct {
	Key string
	Val interface{}
}

type AggregateRelationshipCriterion struct {
	RelationshipCriterion RelationshipCriterion
	Aggregation           db.Aggregation
}

type RelationshipCriterion struct {
	Relationship db.Relationship
	Where        Where
}

type Where struct {
	FieldCriteria                 []FieldCriterion
	RelationshipCriteria          []RelationshipCriterion
	AggregateRelationshipCriteria []AggregateRelationshipCriterion
	Or                            []Where
	And                           []Where
	Not                           []Where
}

func (fc FieldCriterion) Matcher() db.Matcher {
	return db.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) ([]*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleFindMany(tx, root, op.FindArgs)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	return qrs, nil
}

func handleFindMany(tx db.Tx, parent db.ModelRef, fm FindArgs) []db.QueryClause {
	clauses := handleWhere(tx, parent, fm.Where)
	return append(clauses, handleIncludes(tx, parent, fm.Include)...)
}

func handleRelationshipWhere(tx db.Tx, parent db.ModelRef, rel db.Relationship, where Where) (clauses []db.QueryClause, child db.ModelRef) {
	child = tx.Ref(rel.Target().ID())
	j := db.LeftJoin(child, parent.Rel(rel))
	clauses = append(clauses, j)
	if rel.Multi() {
		a := db.Aggregate(child, db.Include)
		clauses = append(clauses, a)
	}
	clauses = append(clauses, handleWhere(tx, child, where)...)
	return clauses, child
}

func getEdgeResults(o, n []*db.QueryResult) []*db.QueryResult {
	if len(o) != len(n) {
		return n
	}
	for i, _ := range o {
		var outs []*db.QueryResult
		if len(o[i].ToOne) != len(n[i].ToOne) {
			for _, v := range n[i].ToOne {
				outs = append(outs, v)
			}
		}
		if len(o[i].ToOne) == len(n[i].ToOne) {
			for k, _ := range n[i].ToOne {
				outs = append(outs, getEdgeResults([]*db.QueryResult{o[i].ToOne[k]}, []*db.QueryResult{n[i].ToOne[k]})...)
			}
		}
		if len(o[i].ToMany) != len(n[i].ToMany) {
			for _, v := range n[i].ToMany {
				outs = append(outs, v...)
			}
		}
		if len(o[i].ToMany) == len(n[i].ToMany) {
			for k, _ := range n[i].ToMany {
				outs = append(outs, getEdgeResults(o[i].ToMany[k], n[i].ToMany[k])...)
			}
		}
		return outs
	}
	return []*db.QueryResult{}
}

func handleWhere(tx db.Tx, parent db.ModelRef, w Where) []db.QueryClause {
	clauses := []db.QueryClause{}
	for _, fc := range w.FieldCriteria {
		clauses = append(clauses, db.Filter(parent, fc.Matcher()))
	}
	for _, rc := range w.RelationshipCriteria {
		clauses = append(clauses, handleRC(tx, parent, rc)...)
	}
	for _, arc := range w.AggregateRelationshipCriteria {
		clauses = append(clauses, handleARC(tx, parent, arc)...)
	}

	var orBlocks []db.QBlock
	for _, or := range w.Or {
		orBlock := handleSetOpBranch(tx, parent, or)
		orBlocks = append(orBlocks, orBlock)
	}
	if len(orBlocks) > 0 {
		clauses = append(clauses, db.Or(parent, orBlocks...))
	}

	var andBlocks []db.QBlock
	for _, and := range w.And {
		andBlock := handleSetOpBranch(tx, parent, and)
		andBlocks = append(andBlocks, andBlock)
	}
	if len(andBlocks) > 0 {
		clauses = append(clauses, db.Union(parent, andBlocks...))
	}

	var notBlocks []db.QBlock
	for _, not := range w.Not {
		notBlock := handleSetOpBranch(tx, parent, not)
		notBlocks = append(notBlocks, notBlock)
	}
	if len(notBlocks) > 0 {
		clauses = append(clauses, db.Not(parent, notBlocks...))
	}
	return clauses
}

func handleSetOpBranch(tx db.Tx, parent db.ModelRef, w Where) db.QBlock {
	clauses := handleWhere(tx, parent, w)
	return db.Subquery(clauses...)
}

func handleRC(tx db.Tx, parent db.ModelRef, rc RelationshipCriterion) []db.QueryClause {
	child := tx.Ref(rc.Relationship.Target().ID())
	on := parent.Rel(rc.Relationship)
	j := db.Join(child, on)
	clauses := handleWhere(tx, child, rc.Where)
	clauses = append(clauses, j)
	return clauses
}

func handleARC(tx db.Tx, parent db.ModelRef, arc AggregateRelationshipCriterion) []db.QueryClause {
	child := tx.Ref(arc.RelationshipCriterion.Relationship.Target().ID())
	on := parent.Rel(arc.RelationshipCriterion.Relationship)

	j := db.Join(child, on)
	a := db.Aggregate(child, arc.Aggregation)
	clauses := handleWhere(tx, child, arc.RelationshipCriterion.Where)
	clauses = append(clauses, a)
	clauses = append(clauses, j)
	return clauses
}
