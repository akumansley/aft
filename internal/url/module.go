package url

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("87770abc-92e0-4cf6-a6a5-8d4c6251213a")
}

func (m *Module) Name() string {
	return "url"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/url"
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		URLValidator,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		URL,
	}
}

var URLValidator = starlark.MakeStarlarkFunction(
	db.MakeID("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	"urlValidator",
	2,
	db.Validator,
	`def main(input, rec):
	# Use a built-in to parse an URL
    u, ok = urlparse(input)
    if not ok:
        # If input is bad, raise an error
        return error(message="Invalid url", code="parse-error")
    return input
`)

var URL = db.MakeCoreDatatype(
	db.MakeID("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	"url",
	db.StringStorage,
	URLValidator,
)
