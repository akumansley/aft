package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/db"
)

type authedDB struct {
	db.DB
}

type authedRWTx struct {
	db.RWTx
}

type authedTx struct {
	db.Tx
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
	return &authedRWTx{RWTx: d.DB.NewRWTxWithContext(ctx)}
}

func (d *authedDB) NewTxWithContext(ctx context.Context) db.Tx {
	return &authedTx{Tx: d.DB.NewTxWithContext(ctx)}
}

func (t *authedTx) WithContext(ctx context.Context) db.Tx {
	return &authedTx{t.Tx.WithContext(ctx)}
}

func (t *authedTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	q := t.Tx.Query(ref, clauses...)
	aq, err := Authed(t.Tx, q, Read)
	if err != nil {
		panic(err)
	}
	return aq
}

func (t *authedRWTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	q := t.RWTx.Query(ref, clauses...)
	aq, err := Authed(t.RWTx, q, Read)
	if err != nil {
		panic(err)
	}
	return aq
}

func (t *authedRWTx) RWWithContext(ctx context.Context) db.RWTx {
	return &authedRWTx{t.RWTx.RWWithContext(ctx)}
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
	case *db.CreateOp:
		create := op.(*db.CreateOp)
		return checkOneRecord(tx, create.Record, Create)
	case *db.DeleteOp:
		return true, nil
	case *db.UpdateOp:
		update := op.(*db.UpdateOp)
		return checkOneRecord(tx, update.NewRecord, Update)
	case *db.ConnectOp:
		connect := op.(*db.ConnectOp)
		return checkConnect(tx, connect.Source, connect.Target, connect.RelID)
	case *db.DisconnectOp:
		disconnect := op.(*db.DisconnectOp)
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
	sourceOK, err := checkOne(tx, source, rel.Source(tx).ID(), Update)
	if err != nil {
		return false, err
	}
	targetOK, err := checkOne(tx, target, rel.Target(tx).ID(), Update)
	if err != nil {
		return false, err
	}
	return sourceOK && targetOK, nil
}

func checkOneRecord(tx db.Tx, rec db.Record, pt PolicyType) (bool, error) {
	return checkOne(tx, rec.ID(), rec.InterfaceID(), pt)
}

func checkOne(tx db.Tx, recID, ifaceID db.ID, pt PolicyType) (bool, error) {
	iface := tx.Ref(ifaceID)
	q := tx.Query(iface, db.Filter(iface, db.EqID(recID)))
	aq, err := Authed(tx, q, pt)
	if err != nil {
		return false, fmt.Errorf("%w error building authed query", err)
	}
	_, err = aq.OneRecord()
	if err == db.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Authed(tx db.Tx, q db.Q, pt PolicyType) (authedQ db.Q, err error) {
	if !shouldAuth(tx.Context()) {
		return q, nil
	}

	var clauses []db.QueryClause

	// filter the root
	if q.Root != nil {
		rootRef := *q.Root
		var c db.QueryClause
		c, err = applyPolicies(tx, rootRef, pt)
		if err != nil {
			return
		}
		clauses = append(clauses, c)
	}
	for _, jl := range q.Joins {
		for _, j := range jl {
			jRef := j.To
			var c db.QueryClause
			c, err = applyPolicies(tx, jRef, pt)
			if err != nil {
				return
			}
			clauses = append(clauses, c)
		}
	}
	for _, sol := range q.SetOps {
		for _, so := range sol {
			var authedBranches []db.Q
			for _, subq := range so.Branches {
				var authQ db.Q
				authQ, err = Authed(tx, subq, pt)
				if err != nil {
					return
				}
				authedBranches = append(authedBranches, authQ)
			}
			so.Branches = authedBranches
		}
	}

	for _, c := range clauses {
		c(&q)
	}
	return q, nil
}

func applyPolicies(tx db.Tx, ref db.ModelRef, pt PolicyType) (db.QueryClause, error) {
	branches := []db.Q{}
	ps := policies(tx, ref.InterfaceID)
	fmt.Printf("policies: %v\n", ps)

	// fail closed
	if len(ps) == 0 {
		return db.Filter(ref, db.False()), nil
	}

	for _, p := range ps {
		clauses, err := applyPolicy(p, tx, ref, pt)
		if err != nil {
			return nil, err
		}
		branches = append(branches, tx.Subquery(clauses...))
	}
	return db.Or(ref, branches...), nil
}

func applyPolicy(p *policy, tx db.Tx, ref db.ModelRef, pt PolicyType) ([]db.QueryClause, error) {
	iface, err := tx.Schema().GetInterfaceByID(ref.InterfaceID)
	if err != nil {
		return nil, err
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
			"$USER_ID": uid,
		}
		subJSON(data, subs)
	}

	w, err := parsers.Parser{tx}.ParseWhere(iface, data)
	if err != nil {
		return nil, err
	}
	clauses := operations.HandleWhere(tx, ref, w)
	return clauses, nil
}

func policies(tx db.Tx, ifaceID db.ID) []*policy {
	role, ok := roleFromContext(tx.Context())
	if !ok {
		panic("No role in context")
	}

	policies := tx.Ref(PolicyModel.ID())
	ifaces := tx.Ref(db.InterfaceInterface.ID())
	roles := tx.Ref(RoleModel.ID())
	policyRole, _ := tx.Schema().GetRelationshipByID(PolicyRole.ID())
	policyFor, _ := tx.Schema().GetRelationshipByID(PolicyFor.ID())

	q := tx.Query(policies,
		db.Join(roles, policies.Rel(policyRole)),
		db.Filter(roles, db.EqID(role.ID())),
		db.Aggregate(roles, db.Some),
		db.Join(ifaces, policies.Rel(policyFor)),
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
