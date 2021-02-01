package db

import "fmt"

// Model

var EnumModel = MakeModel(
	MakeID("7f0d30ae-83c7-4d89-a134-e8ac326321e6"),
	"enum",
	[]AttributeL{
		MakeConcreteAttribute(
			MakeID("57ff4302-41b4-4866-8a5c-bc1c264ffba4"),
			"name",
			String,
		),
		eSystem,
	},
	[]RelationshipL{
		EnumEnumValues,
	},
	[]ConcreteInterfaceL{DatatypeInterface},
)

var EnumEnumValues = MakeConcreteRelationship(
	MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
	"enumValues",
	true,
	EnumValueModel,
)

var eSystem = MakeConcreteAttribute(
	MakeID("5b3838f4-fa33-4723-9ea7-eed5928220fd"),
	"system",
	Bool,
)

// Loader

type EnumDatatypeLoader struct{}

func (l EnumDatatypeLoader) ProvideModel() ModelL {
	return EnumModel
}

func (l EnumDatatypeLoader) Load(rec Record) Datatype {
	return &enum{rec}
}

// Literal

func MakeEnum(id ID, name string, values []EnumValueL) EnumL {
	return EnumL{
		id,
		true,
		name,
		values,
	}
}

type EnumL struct {
	ID_     ID     `record:"id"`
	System  bool   `record:"system"`
	Name_   string `record:"name"`
	Values_ []EnumValueL
}

func (lit EnumL) ID() ID {
	return lit.ID_
}

func (lit EnumL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)

	recs = append(recs, rec)
	for _, a := range lit.Values_ {
		ars, al := a.MarshalDB(b)
		recs = append(recs, ars...)
		links = append(links, al...)

		links = append(links, Link{rec.ID(), a.ID(), EnumEnumValues})
	}

	return
}

func (lit EnumL) InterfaceID() ID {
	return EnumModel.ID()
}

func (lit EnumL) InterfaceName() string {
	return EnumModel.Name_
}

func (lit EnumL) Load(tx Tx) Datatype {
	dt, err := tx.Schema().GetDatatypeByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return dt
}

// Dynamic

type enum struct {
	rec Record
}

func (e *enum) ID() ID {
	return e.rec.ID()
}

func (e *enum) Name() string {
	return e.rec.MustGet("name").(string)
}

func (e *enum) Storage(tx Tx) EnumValue {
	ev, err := tx.Schema().GetEnumValueByID(UUIDStorage.ID())
	if err != nil {
		err := fmt.Errorf("UUIDStorage error %w", err)
		panic(err)
	}
	return ev
}

func (e *enum) FromJSON(tx Tx) (Function, error) {
	f, err := tx.Schema().GetFunctionByID(uuidValidator.ID())
	return f, err
}
