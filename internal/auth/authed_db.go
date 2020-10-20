package auth

import (
	"context"

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
	return Authed(t.Tx, q, t.ctx)
}

func (t *authedRWTx) Query(ref db.ModelRef, clauses ...db.QueryClause) db.Q {
	q := t.RWTx.Query(ref, clauses...)
	return Authed(t.RWTx, q, t.ctx)
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

func Authed(tx db.Tx, q db.Q, ctx context.Context) db.Q {
	if ctx == noAuthContext {
		return q
	}

	userRec, ok := FromContext(tx, ctx)
	user := &user{userRec, tx}
	var role db.Record
	var err error
	if !ok {
		roles := tx.Ref(RoleModel.ID())
		role, err = tx.Query(roles, db.Filter(roles, db.EqID(Public.ID()))).OneRecord()
		if err != nil {
			panic("Couldn't find public role")
		}
	} else {
		role, err = getRole(tx, user)
		if err != nil {
			panic(err)
		}
	}
	var clauses []db.QueryClause

	// filter the root
	if q.Root != nil {
		rootRef := *q.Root
		ps := policies(tx, role, rootRef.InterfaceID)
		clauses = append(clauses, applyPolicies(tx, ps, rootRef, user))
	}
	for _, jl := range q.Joins {
		for _, j := range jl {
			jRef := j.To
			ps := policies(tx, role, jRef.InterfaceID)
			clauses = append(clauses, applyPolicies(tx, ps, jRef, user))
		}
	}
	for _, sol := range q.SetOps {
		for _, so := range sol {
			var authedBranches []db.Q
			for _, subq := range so.Branches {
				authedBranches = append(authedBranches, Authed(tx, subq, ctx))
			}
			so.Branches = authedBranches
		}
	}

	for _, c := range clauses {
		c(&q)
	}
	return q
}

func applyPolicies(tx db.Tx, ps []*policy, ref db.ModelRef, user *user) db.QueryClause {
	branches := []db.Q{}

	// fail closed
	if len(ps) == 0 {
		return db.Filter(ref, db.False())
	}

	for _, p := range ps {
		clauses := p.Apply(tx, ref, user)
		branches = append(branches, db.Subquery(clauses...))
	}
	return db.Or(ref, branches...)
}

func getRole(tx db.Tx, user *user) (db.Record, error) {
	roles := tx.Ref(RoleModel.ID())
	users := tx.Ref(UserModel.ID())

	q := tx.Query(roles,
		db.Join(users, roles.Rel(RoleUsers)),
		db.Aggregate(users, db.Some),
		db.Filter(users, db.EqID(user.ID())),
	)
	roleRec, err := q.OneRecord()

	return roleRec, err
}

func policies(tx db.Tx, role db.Record, ifaceID db.ID) []*policy {
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
