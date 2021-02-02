package auth

import (
	"context"
	"errors"
	"net/http"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

var noAuthContext = context.WithValue(context.Background(), noAuthKey, true)

type key int

const (
	userKey key = iota
	roleKey
	noAuthKey
	setCookieKey
)

func (k key) String() string {
	switch k {
	case userKey:
		return "user"
	case roleKey:
		return "role"
	case noAuthKey:
		return "noAuth"
	case setCookieKey:
		return "setCookie"
	}
	panic("invalid key")
}

func WithRole(ctx context.Context, role db.Record) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func RoleFromContext(ctx context.Context) (db.Record, bool) {
	r, ok := ctx.Value(roleKey).(db.Record)
	return r, ok
}

func WithUser(ctx context.Context, user db.Record) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (db.Record, bool) {
	u, ok := ctx.Value(userKey).(db.Record)
	return u, ok
}

func IDFromContext(ctx context.Context) (id db.ID, ok bool) {
	u, ok := ctx.Value(userKey).(db.Record)
	if ok {
		id = u.ID()
	}
	return
}

func withSetCookie(ctx context.Context, w http.ResponseWriter) context.Context {
	setCookie := func(cookie *http.Cookie) {
		http.SetCookie(w, cookie)
	}
	return context.WithValue(ctx, setCookieKey, setCookie)
}

func setCookieFromContext(ctx context.Context) (setCookie func(*http.Cookie), ok bool) {
	v, ok := ctx.Value(setCookieKey).(func(*http.Cookie))
	return v, ok
}

func makeAuthMiddleware(appDB db.DB) lib.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// add setCookie no matter what
			ctx := withSetCookie(r.Context(), w)

			tx := appDB.NewTxWithContext(noAuthContext)
			tokCookie, err := r.Cookie("tok")
			if err == nil {
				token := tokCookie.Value
				user, err := UserForToken(tx, token)

				if err == nil {
					ctx = WithUser(ctx, user)
					role, err := getRole(tx, user)
					if err == db.ErrNotFound {
						role = getPublic(tx)
					} else if err != nil {
						lib.WriteError(w, err)
						return
					}

					ctx = WithRole(ctx, role)
				} else if !errors.Is(err, ErrInvalid) {
					lib.WriteError(w, err)
					return
				}
			} else {
				role := getPublic(tx)
				ctx = WithRole(ctx, role)
			}
			r = r.Clone(ctx)

			inner.ServeHTTP(w, r)
		})
	}
}

func getRole(tx db.Tx, user db.Record) (db.Record, error) {
	roles := tx.Ref(RoleModel.ID())
	users := tx.Ref(UserModel.ID())
	roleUsers, err := tx.Schema().GetRelationshipByID(RoleUsers.ID())
	if err != nil {
		return nil, err
	}

	q := tx.Query(roles,
		db.Join(users, roles.Rel(roleUsers)),
		db.Aggregate(users, db.Some),
		db.Filter(users, db.EqID(user.ID())),
	)
	roleRec, err := q.OneRecord()

	return roleRec, err
}
