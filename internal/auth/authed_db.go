package auth

import (
	"awans.org/aft/internal/db"
	"context"
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
	return d.DB.NewRWTx()
}

func (d *authedDB) NewTx() db.Tx {
	return d.DB.NewTx()
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

func Authed(tx db.Tx, q db.Q, ctx context.Context) db.Q {
	user, ok := FromContext(tx, ctx)
	if !ok {
		// TODO just replace this with a "public" role
		panic("No user")
	}
	var clauses []db.QueryClause

	// filter the root
	if q.Root != nil {
		rootRef := *q.Root
		ps := policies(tx, user, rootRef.InterfaceID)
		clauses = append(clauses, applyPolicies(tx, ps, rootRef))
	}
	for _, jl := range q.Joins {
		for _, j := range jl {
			jRef := j.To
			ps := policies(tx, user, jRef.InterfaceID)
			clauses = append(clauses, applyPolicies(tx, ps, jRef))
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

func applyPolicies(tx db.Tx, ps []*policy, ref db.ModelRef) db.QueryClause {
	branches := []db.Q{}
	for _, p := range ps {
		clauses := p.Apply(tx, ref)
		branches = append(branches, db.Subquery(clauses...))
	}
	return db.Or(ref, branches...)
}

func policies(tx db.Tx, user user, modelID db.ID) []*policy {
	policies := tx.Ref(PolicyModel.ID())
	models := tx.Ref(db.ModelModel.ID())
	roles := tx.Ref(RoleModel.ID())
	users := tx.Ref(UserModel.ID())

	q := tx.Query(policies,
		db.Join(roles, policies.Rel(PolicyRoles)),
		db.Aggregate(roles, db.Some),

		db.Join(users, roles.Rel(RoleUsers)),
		db.Aggregate(users, db.Some),
		db.Filter(users, db.EqID(user.ID())),

		db.Join(models, policies.Rel(PolicyFor)),
		db.Filter(models, db.EqID(modelID)),
	)
	results := q.Records()

	constructed := []*policy{}
	for _, p := range results {
		constructed = append(constructed, &policy{p, tx})
	}

	return constructed
}
