package db

type nodeType int

func Plan(q Q) Node {
	return PlanRoot(q, *q.Root)
}

func PlanJoinOuter(q Q, j JoinOperation) (n Node) {
	ref := j.To
	rel := j.on.rel
	aliasID := ref.AliasID
	order, _ := q.Orderings[aliasID]
	projection, _ := q.Selections[aliasID]

	iface, _ := q.tx.Schema().GetInterfaceByID(ref.InterfaceID)
	n = &RelLookupNode{
		tx:          q.tx,
		interfaceID: ref.InterfaceID,
		rel:         rel,
		iface:       iface,
		order:       order,
		projection:  projection,
	}

	n = PlanRef(q, n, ref)
	return n
}

func PlanCaseOuter(q Q, inner Node, c CaseOperation) (n Node) {
	n = inner
	ref := c.Of
	n = PlanRef(q, n, ref)
	return n
}

func PlanJoins(q Q, inner Node, innerRef ModelRef) (n Node) {
	n = inner

	aliasID := innerRef.AliasID
	for _, join := range q.Joins[aliasID] {
		outer := PlanJoinOuter(q, join)
		aggregation := q.Aggregations[join.To.AliasID]
		filters, _ := q.Filters[join.To.AliasID]
		n = &JoinNode{
			inner:       n,
			outer:       outer,
			rel:         join.on.rel,
			filters:     filters,
			aggregation: aggregation,
			joinType:    join.jt,
		}
	}
	return n
}

func PlanCases(q Q, inner Node, innerRef ModelRef) (n Node) {
	n = inner

	var caseNodes []Node
	aliasID := innerRef.AliasID
	var notFilters []Matcher
	for _, c := range q.Cases[aliasID] {
		isCase := &FilterNode{
			inner:   n,
			matcher: IsModel(c.Of.InterfaceID),
		}
		notFilters = append(notFilters, IsNotModel(c.Of.InterfaceID))
		caseWithOuter := PlanCaseOuter(q, isCase, c)
		caseNodes = append(caseNodes, caseWithOuter)
	}

	if len(caseNodes) > 0 {
		isNotCase := &FilterNode{
			inner:   n,
			matcher: And(notFilters...),
		}

		n = &UnionNode{
			nodes: append(caseNodes, isNotCase),
		}
	}

	return
}

func PlanSetOps(q Q, inner Node, ref ModelRef) (n Node) {
	n = inner

	aliasID := ref.AliasID
	for _, setOp := range q.SetOps[aliasID] {
		switch setOp.op {
		case or:
			branchNodes := []Node{}
			for _, branch := range setOp.Branches {
				branchNode := PlanRef(branch, n, ref)
				branchNodes = append(branchNodes, branchNode)
			}
			n = &UnionNode{
				nodes: branchNodes,
			}
		case not:
			for _, branch := range setOp.Branches {

				// not sure about the PlanRef
				branchNode := PlanRef(branch, n, ref)
				n = &SubtractNode{
					left:  n,
					right: branchNode,
				}
			}
		case and:
			for _, branch := range setOp.Branches {
				// not sure about the PlanRef
				branchNode := PlanRef(branch, n, ref)
				n = &IntersectionNode{
					left:  n,
					right: branchNode,
				}
			}
		}
	}
	return
}

func PlanRoot(q Q, ref ModelRef) (n Node) {
	aliasID := ref.AliasID
	var filters []Matcher
	filters, _ = q.Filters[aliasID]

	var order []Sort
	order, _ = q.Orderings[aliasID]

	var projection Selection
	projection, _ = q.Selections[aliasID]

	// the leftmost is always a table scan
	// that "drives" all of the other joins
	iface, _ := q.tx.Schema().GetInterfaceByID(ref.InterfaceID)
	n = &TableAccessNode{
		interfaceID: ref.InterfaceID,
		iface:       iface,
		filters:     filters,
		order:       order,
		projection:  projection,
	}
	n = PlanRef(q, n, ref)

	return n
}

func PlanRef(q Q, inner Node, ref ModelRef) (n Node) {
	aliasID := ref.AliasID
	n = inner
	filters, ok := q.Filters[aliasID]
	if ok {
		n = &FilterNode{
			inner:   n,
			matcher: And(filters...),
		}
	}
	n = PlanCases(q, n, ref)
	n = PlanJoins(q, n, ref)
	n = PlanSetOps(q, n, ref)

	// TODO: should these come before or after setOp?
	if offset, ok := q.Offsets[aliasID]; ok && offset != 0 {
		n = &OffsetNode{
			offset: offset,
			inner:  n,
		}
	}

	if limit, ok := q.Limits[aliasID]; ok && limit != 0 {
		n = &LimitNode{
			limit: limit,
			inner: n,
		}
	}
	return
}
