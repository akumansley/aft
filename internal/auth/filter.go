package auth

import (
	"context"
	"fmt"
	"net/http"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

var noAuthContext = context.WithValue(context.Background(), noAuthKey, true)

var userKey = "user"
var noAuthKey = "noAuth"

func WithUser(ctx context.Context, user db.Record) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func FromContext(tx db.Tx, ctx context.Context) (*user, bool) {
	u, ok := ctx.Value(userKey).(db.Record)
	if ok {
		return &user{u, tx}, ok
	}
	return nil, false
}

func makeAuthMiddleware(appDB db.DB) lib.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "" {
				user, err := UserForToken(appDB, token)
				fmt.Printf("AUTH: user: %v err: %v\n", user, err)
				ctx := WithUser(r.Context(), user)
				r = r.Clone(ctx)
			}

			inner.ServeHTTP(w, r)
		})
	}
}
