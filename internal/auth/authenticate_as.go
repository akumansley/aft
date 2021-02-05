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
	0,
	db.Internal,
	clearAuthentication,
)

func clearAuthentication(ctx context.Context, args []interface{}) (result interface{}, err error) {
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

func authenticateAsFunc(ctx context.Context, args []interface{}) (result interface{}, err error) {
	id := args[0].(uuid.UUID)
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		return nil, errors.New("No tx found in authenticateAs")
	}
	ActAsFunction(rwtx)
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
	ActAsUser(rwtx)
	return
}

var AuthenticateAs = MakeNativeFunctionWithRole(
	db.MakeID("e20ae44f-6a5e-4d25-ab13-de3bd7b7c392"),
	"authenticateAs",
	1,
	db.Internal,
	authenticateAsFunc,
	LoginSystem,
)
