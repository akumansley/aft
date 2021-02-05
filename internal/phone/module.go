package phone

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("61c603a9-224e-4da7-b6a7-574e4ee2a0e5")
}

func (m *Module) Name() string {
	return "phone"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/phone"
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		PhoneValidator,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		Phone,
	}
}

var PhoneValidator = starlark.MakeStarlarkFunction(
	db.MakeID("f720efdc-3694-429f-9d4e-c2150388bd30"),
	"phone",
	2,
	db.Validator,
	`# Compile Regular Expression for valid US Phone Numbers
phone = re.compile(r"^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$")

def main(input, rec):
    if not phone.match(input):
        return error(message="Invalid phone number", code="parse-error")
    # Otherwise, return it stripped of formatting
    clean = input.replace(" ","").replace("-","")
    return clean.replace("(","").replace(")","")`,
)

var Phone = db.MakeCoreDatatype(
	db.MakeID("d5b7bc19-9eec-4bf9-b362-1a642458060f"),
	"phone",
	db.IntStorage,
	PhoneValidator,
)
