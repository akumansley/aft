package auth

import (
	"errors"
	"net/http"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

func makeAuthMiddleware(appDB db.DB) lib.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// add setCookie no matter what
			ctx := withSetCookie(r.Context(), w)

			tx := appDB.NewTxWithContext(r.Context())
			deescalate := Escalate(tx)

			tokCookie, err := r.Cookie("tok")
			if err == nil {
				token := tokCookie.Value
				user, err := UserForToken(tx, token)

				if err == nil {
					ctx = withUser(ctx, user)
					role, err := RoleForUser(tx, user)
					if err == db.ErrNotFound {
						role = getPublic(tx)
					} else if err != nil {
						lib.WriteError(w, err)
						return
					}

					ctx = withRole(ctx, role)
				} else if !errors.Is(err, ErrInvalid) {
					lib.WriteError(w, err)
					return
				}
			} else {
				role := getPublic(tx)
				ctx = withRole(ctx, role)
			}
			r = r.Clone(ctx)

			deescalate()
			inner.ServeHTTP(w, r)
		})
	}
}

func getPublic(tx db.Tx) db.Record {
	roles := tx.Ref(RoleModel.ID())
	val, err := tx.Query(roles, db.Filter(roles, db.EqID(Public.ID()))).OneRecord()
	if err != nil {
		panic(err)
	}
	return val
}
