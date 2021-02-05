package auth

import (
	"context"
	"fmt"
	"regexp"

	"awans.org/aft/internal/db"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func emailFromJSON(ctx context.Context, args []interface{}) (interface{}, error) {
	value := args[0]
	email, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	valid := isEmailValid(email)
	if !valid {
		return nil, fmt.Errorf("%w: expected valid email got %v", ErrValue, value)
	}

	return email, nil
}

var emailAddressValidator = db.MakeNativeFunction(
	db.MakeID("ed046b08-ade2-4570-ade4-dd1e31078219"),
	"emailAddress",
	1,
	db.Validator,
	emailFromJSON)

var EmailAddress = db.MakeCoreDatatype(
	db.MakeID("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	"emailAddress",
	db.StringStorage,
	emailAddressValidator,
)
