package auth

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

var (
	ErrAccount = fmt.Errorf("%w: unable to create account", ErrAuth)
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Data interface{} `json:"data"`
}

type SignupHandler struct {
	bus *bus.EventBus
	db  db.DB
}

func (sh SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	var sr SignupRequest
	buf, _ := ioutil.ReadAll(r.Body)
	err = jsoniter.Unmarshal(buf, &sr)
	if err != nil {
		return
	}

	sh.bus.Publish(lib.ParseRequest{Request: sr})

	rwtx := sh.db.NewRWTx()
	user, err := rwtx.FindOne(UserModel.ID(), db.Eq("email", sr.Email))
	if !errors.Is(err, db.ErrNotFound) {
		return ErrAccount
	}
	user = db.RecordForModel(UserModel)
	uid := uuid.New()
	user.Set("id", uid)
	user.Set("email", sr.Email)
	user.Set("password", sr.Password)
	rwtx.Insert(user)
	rwtx.Commit()

	response := SignupResponse{Data: user}

	// write out the response
	bytes, _ := jsoniter.Marshal(&response)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
	return nil
}
