package repl

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"fmt"
	"strings"
)

func eval(input string, tx db.RWTx) string {
	sh := starlark.StarlarkFunctionHandle{Code: input, Env: starlark.ReplLib(tx)}
	r, err := sh.Invoke("")
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
