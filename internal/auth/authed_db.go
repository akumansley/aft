package auth

import (
	"context"
	"encoding/json"
	"errors"

	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/db"
)

type authedDB struct {
	db.DB
}

type authedRWTx struct {
	db.RWTx
	ctx context.Context
}

type authedTx struct {
	db.Tx
	ctx context.Context
}

func AuthedDB(d db.DB) db.DB {
	return &authedDB{DB: d}
}

func (d *authedDB) NewRWTx() db.RWTx {
	return d.NewRWTxWithContext(context.Background())
}

func (d *authedDB) NewTx() db.Tx {
	return d.NewTxWithContext(context.Background())
}

func (d *authedDB) NewRWTxWithContext(ctx context.Context) db.RWTx {
	return &authedRWTx{RWTx: d.DB.NewRWTxWithContext(ctx), ctx: ctx}
}

func (d *authedDB) NewTxWithContext(ctx context.Context) db.Tx {
	return &authedTx{Tx: d.DB.NewTxWithContext(ctx), ctx: ctx}
}

func (t *authedTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	q := t.Tx.Query(ref, clauses...)
	return Authed(t.Tx, q, Read)
}

func (t *authedRWTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	q := t.RWTx.Query(ref, clauses...)
	return Authed(t.RWTx, q, Read)
}

func (t *authedRWTx) Schema() *db.Schema {
	s := t.RWTx.Schema()
	s.SetTx(t)
	return s
}

func (t *authedTx) Schema() *db.Schema {
	s := t.Tx.Schema()
	s.SetTx(t)
	return s
}

func (t *authedRWTx) Insert(rec db.Record) error {
	return t.RWTx.Insert(rec)
}

func (t *authedRWTx) Update(oldRec, newRec db.Record) error {
	return t.RWTx.Update(oldRec, newRec)
}

func (t *authedRWTx) Connect(sourceID, targetID, relID db.ID) error {
	return t.RWTx.Connect(sourceID, targetID, relID)
}

func (t *authedRWTx) Disconnect(sourceID, targetID, relID db.ID) error {
	return t.RWTx.Disconnect(sourceID, targetID, relID)
}

func PostconditionHandler(event db.BeforeCommit) {
	ops := event.Tx.Operations()
	if len(ops) == 0 {
		return
	}
	for _, op := range ops {
		authorized, err := checkOpPostcondition(event.Tx, op)
		if err != nil || !authorized {
			event.Tx.Abort(ErrAuth)
			return
		}
	}
}

func checkOpPostcondition(tx db.Tx, op db.Operation) (bool, error) {
	switch op.(type) {
	case db.CreateOp:
		create := op.(db.CreateOp)
		return checkOneRecord(tx, create.Record, Create)
	case db.DeleteOp:
		return true, nil
	case db.UpdateOp:
		update := op.(db.UpdateOp)
		return checkOneRecord(tx, update.NewRecord, Update)
	case db.ConnectOp:
		connect := op.(db.ConnectOp)
		return checkConnect(tx, connect.Source, connect.Target, connect.RelID)
	case db.DisconnectOp:
		disconnect := op.(db.DisconnectOp)
		return checkConnect(tx, disconnect.Source, disconnect.Target, disconnect.RelID)
	default:
		return false, errors.New("Invalid op")
	}
}

func checkConnect(tx db.Tx, source, target, relID db.ID) (bool, error) {
	rel, err := tx.Schema().GetRelationshipByID(relID)
	if err != nil {
		return false, err
	}
	sourceOK, err := checkOne(tx, source, rel.Source().ID(), Update)
	if err != nil {
		return false, err
	}
	targetOK, err := checkOne(tx, target, rel.Target().ID(), Update)
	if err != nil {
		return false, err
	}
	return sourceOK && targetOK, nil
}

func checkOneRecord(tx db.Tx, rec db.Record, pt PolicyType) (bool, error) {
	return checkOne(tx, rec.ID(), rec.Interface().ID(), pt)
}

func checkOne(tx db.Tx, recID, ifaceID db.ID, pt PolicyType) (bool, error) {
	iface := tx.Ref(ifaceID)
	q := tx.Query(iface, db.Filter(iface, db.EqID(recID)))
	aq := Authed(tx, q, pt)
	_, err := aq.OneRecord()
	if err == db.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Authed(tx db.Tx, q db.Q, pt PolicyType) db.Q {
	if tx.Context() == noAuthContext {
		return q
	}

	var clauses []db.QueryClause

	// filter the root
	if q.Root != nil {
		rootRef := *q.Root
		clauses = append(clauses, applyPolicies(tx, rootRef, pt))
	}
	for _, jl := range q.Joins {
		for _, j := range jl {
			jRef := j.To
			clauses = append(clauses, applyPolicies(tx, jRef, pt))
		}
	}
	for _, sol := range q.SetOps {
		for _, so := range sol {
			var authedBranches []db.Q
			for _, subq := range so.Branches {
				authedBranches = append(authedBranches, Authed(tx, subq, pt))
			}
			so.Branches = authedBranches
		}
	}

	for _, c := range clauses {
		c(&q)
	}
	return q
}

func applyPolicies(tx db.Tx, ref db.ModelRef, pt PolicyType) db.QueryClause {
	branches := []db.Q{}
	ps := policies(tx, ref.InterfaceID)

	// fail closed
	if len(ps) == 0 {
		return db.Filter(ref, db.False())
	}

	for _, p := range ps {
		clauses := applyPolicy(p, tx, ref, pt)
		branches = append(branches, db.Subquery(clauses...))
	}
	return db.Or(ref, branches...)
}

func applyPolicy(p *policy, tx db.Tx, ref db.ModelRef, pt PolicyType) []db.QueryClause {
	iface, err := tx.Schema().GetInterfaceByID(ref.InterfaceID)
	if err != nil {
		panic("bad")
	}

	// TODO check allowRead etc also
	var templateText string
	switch pt {
	case Read:
		templateText = p.ReadWhere()
	case Create:
		templateText = p.CreateWhere()
	case Update:
		templateText = p.UpdateWhere()
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(templateText), &data)

	user, ok := UserFromContext(tx.Context())
	if ok {
		uid := user.ID().String()
		subs := map[string]interface{}{
			"$userID": uid,
		}
		subJSON(data, subs)
	}

	w, err := parsers.Parser{tx}.ParseWhere(iface, data)
	if err != nil {
		panic(err)
	}
	clauses := operations.HandleWhere(tx, ref, w)
	return clauses
}

func policies(tx db.Tx, ifaceID db.ID) []*policy {
	role, ok := RoleFromContext(tx.Context())
	if !ok {
		panic("No role in context")
	}

	policies := tx.Ref(PolicyModel.ID())
	ifaces := tx.Ref(db.InterfaceInterface.ID())
	roles := tx.Ref(RoleModel.ID())

	q := tx.Query(policies,
		db.Join(roles, policies.Rel(PolicyRole)),
		db.Filter(roles, db.EqID(role.ID())),
		db.Aggregate(roles, db.Some),
		db.Join(ifaces, policies.Rel(PolicyFor)),
		db.Filter(ifaces, db.EqID(ifaceID)),
		db.Aggregate(ifaces, db.Some),
	)
	results := q.Records()

	constructed := []*policy{}
	for _, p := range results {
		constructed = append(constructed, &policy{p, tx})
	}
	return constructed
}
