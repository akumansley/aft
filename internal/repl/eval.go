package repl

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"fmt"
	"strings"
)

func eval(input string, tx db.RWTx) string {
	sf := starlark.MakeStarlarkFunction(db.NewID(), "", db.RPC, input)
	r, err := sf.Call("")
	// An error in eval isn't an error, so make it the result
	if err != nil {
		errString := fmt.Sprintf("%s", err)
		errString = strings.TrimLeft(errString, ":")
		r = fmt.Sprintf("Starlark: %s", errString)
	}
	if r == nil {
		r = ""
	}
	return fmt.Sprintf("%v", r)
}
