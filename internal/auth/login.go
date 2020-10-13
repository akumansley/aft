package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	jsoniter "github.com/json-iterator/go"
)

var (
	ErrAuth = errors.New("auth-error")
	// this is generic so we don't return whether an email has an acct
	ErrUnsuccessful = fmt.Errorf("%w: login unsuccessful", ErrAuth)
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Data interface{} `json:"data"`
}

type LoginHandler struct {
	bus *bus.EventBus
	db  db.DB
}

var ttl = 30 * time.Minute

func authenticateAsFunc(args []interface{}) (result interface{}, err error) {
	ctx := args[0].(context.Context)
	id := args[1].(db.ID)
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		return nil, errors.New("No tx found in authenticateAs")
	}

	tok, err := TokenForID(rwtx, id)
	if err != nil {
		return
	}

	setCookie, ok := setCookieFromContext(ctx)
	if !ok {
		return nil, errors.New("No setCookie found in authenticateAs")
	}

	expires := time.Now().Add(ttl)

	cookie := http.Cookie{
		Name:    "tok",
		Value:   tok,
		Expires: expires,
		Domain:  "",
		Path:    "/",
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

func getUserByEmail(tx db.Tx, email string) (*user, error) {
	users := tx.Ref(UserModel.ID())
	userRec, err := tx.Query(users, db.Filter(users, db.Eq("email", email))).OneRecord()
	if err != nil {
		return nil, err
	}
	return &user{userRec, tx}, nil

}

func writeLogin(w http.ResponseWriter, user *user) {
	response := LoginResponse{Data: user.rec}
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}
