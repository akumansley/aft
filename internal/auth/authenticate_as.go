package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var ttl = 30 * time.Minute

var ClearAuthentication = db.MakeNativeFunction(
	db.MakeID("13ed6cf8-94a0-4732-a673-9c0a0e9c656f"),
	"clearAuthentication",
	1,
	clearAuthentication,
)

func clearAuthentication(args []interface{}) (result interface{}, err error) {
	ctx := args[0].(context.Context)
	setCookie, ok := setCookieFromContext(ctx)
	if !ok {
		return nil, errors.New("No setCookie found in clearAuthentication")
	}

	expires := time.Now()

	cookie := http.Cookie{
		Name:     "tok",
		Value:    "",
		Expires:  expires,
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	setCookie(&cookie)
	return
}

func authenticateAsFunc(args []interface{}) (result interface{}, err error) {
	ctx := args[0].(context.Context)
	id := args[1].(uuid.UUID)
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		return nil, errors.New("No tx found in authenticateAs")
	}

	tok, err := TokenForID(rwtx, db.ID(id))
	if err != nil {
		return
	}

	setCookie, ok := setCookieFromContext(ctx)
	if !ok {
		return nil, errors.New("No setCookie found in authenticateAs")
	}

	expires := time.Now().Add(ttl)

	cookie := http.Cookie{
		Name:     "tok",
		Value:    tok,
		Expires:  expires,
		Domain:   "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	setCookie(&cookie)
	return
}

var AuthenticateAs = db.MakeNativeFunction(
	db.MakeID("e20ae44f-6a5e-4d25-ab13-de3bd7b7c392"),
	"authenticateAs",
	2,
	authenticateAsFunc,
)
