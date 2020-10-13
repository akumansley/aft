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
	noAuthKey
	setCookieKey
)

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
			r = r.Clone(ctx)

			token := r.Header.Get("Authorization")
			if token != "" {
				tx := appDB.NewTxWithContext(noAuthContext)
				user, err := UserForToken(tx, token)

				if err == nil {
					ctx := WithUser(r.Context(), user)
					r = r.Clone(ctx)
				} else if !errors.Is(err, ErrInvalid) {
					lib.WriteError(w, err)
					return
				}

			}

			inner.ServeHTTP(w, r)
		})
	}
}
