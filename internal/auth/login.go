package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	jsoniter "github.com/json-iterator/go"
	"github.com/markbates/pkger"
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

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	tx := lh.db.NewTxWithContext(noAuthContext)

	ctx := r.Context()
	user, ok := FromContext(tx, ctx)

	if ok {
		writeLogin(w, user)
		return
	}
	f, err := pkger.Open("/internal/auth/login.star")
	defer f.Close()

	b2, err := ioutil.ReadAll(f)

	fmt.Printf("code: %v /code\n err: %v\n", string(b2), err)
	if err != nil {
		return err
	}

	var lr LoginRequest
	buf, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(buf, &lr)
	if err != nil {
		return
	}

	lh.bus.Publish(lib.ParseRequest{Request: lr})

	user, err = getUserByEmail(tx, lr.Email)

	if err != nil {
		return ErrUnsuccessful
	}

	pw := user.Password()

	if lr.Password != pw {
		return ErrUnsuccessful
	}

	tok, err := TokenForUser(lh.db, user)
	if err != nil {
		return
	}

	expires := time.Now().Add(ttl)

	cookie := http.Cookie{
		Name:    "tok",
		Value:   tok,
		Expires: expires,
		Domain:  "",
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	writeLogin(w, user)

	return
}

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
