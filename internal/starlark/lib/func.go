package lib

import (
	"context"
	"errors"
	"fmt"

	"awans.org/aft/internal/db"
	"go.starlark.net/starlark"
)

type magicFunc struct {
	ctx context.Context
}

func (c *magicFunc) String() string {
	return "func()"
}

func (c *magicFunc) Type() string {
	return "func"
}

func (c *magicFunc) Freeze() {}

func (c *magicFunc) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (c *magicFunc) Hash() (uint32, error) {
	return 0, errors.New("Unhashable")
}

// returns (nil, nil) if attribute not present
func (c *magicFunc) Attr(name string) (starlark.Value, error) {
	tx, ok := db.TxFromContext(c.ctx)
	if !ok {
		return starlark.None, fmt.Errorf("No tx in contxt")
	}
	f, err := tx.Schema().GetFunctionByName(name)
	if err != nil {
		return starlark.None, fmt.Errorf("No function %v - %v", name, err)
	}
	return makeBuiltin(name, f.Arity()), nil

}

func (c *magicFunc) AttrNames() (names []string) {
	tx, ok := db.TxFromContext(c.ctx)
	if !ok {
		return []string{}
	}
	function := tx.Ref(db.FunctionInterface.ID())
	recs := tx.Query(function).Records()
	for _, r := range recs {
		name := r.MustGet("name").(string)
		names = append(names, name)
	}

	return
}
