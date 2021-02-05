package auth

import (
	"context"
	"net/http"

	"awans.org/aft/internal/db"
)

type key int

const (
	userKey key = iota
	userRoleKey
	roleKey
	functionRoleKey
	noAuthKey
	setCookieKey
)

func (k key) String() string {
	switch k {
	case userKey:
		return "user"
	case userRoleKey:
		return "userRole"
	case roleKey:
		return "role"
	case noAuthKey:
		return "noAuth"
	case setCookieKey:
		return "setCookie"
	case functionRoleKey:
		return "functionRole"
	}
	panic("invalid key")
}

func UserFromContext(ctx context.Context) (db.Record, bool) {
	u, ok := ctx.Value(userKey).(db.Record)
	return u, ok
}

func withUser(ctx context.Context, user db.Record) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func withRole(ctx context.Context, role db.Record) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func roleFromContext(ctx context.Context) (db.Record, bool) {
	r, ok := ctx.Value(roleKey).(db.Record)
	return r, ok
}

func withUserRole(ctx context.Context, role db.Record) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func withFunctionRole(ctx context.Context, role db.Record) context.Context {
	return context.WithValue(ctx, functionRoleKey, role)
}

func functionRoleFromContext(ctx context.Context) (db.Record, bool) {
	u, ok := ctx.Value(functionRoleKey).(db.Record)
	return u, ok
}

func userRoleFromContext(ctx context.Context) (db.Record, bool) {
	u, ok := ctx.Value(userRoleKey).(db.Record)
	return u, ok
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
