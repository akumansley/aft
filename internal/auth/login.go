package auth

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
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

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	var lr LoginRequest
	buf, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(buf, &lr)
	if err != nil {
		return
	}

	lh.bus.Publish(lib.ParseRequest{Request: lr})
	tx := lh.db.NewTx()
	user, err := tx.FindOne(UserModel.Name, db.Eq("email", lr.Email))
	if err != nil {
		return ErrUnsuccessful
	}

	if lr.Password != user.Get("password") {
		return ErrUnsuccessful
	}

	response := LoginResponse{Data: user}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return
}
