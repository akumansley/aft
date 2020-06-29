package db

// Model

var EnumValueModel = MakeModel(
	MakeID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"),
	"enumValue",
	[]AttributeL{
		evName,
	},
	[]RelationshipL{},
	[]ConcreteInterfaceL{},
)

var evName = MakeConcreteAttribute(
	MakeID("5803e350-48f8-448d-9901-7c80f45c775b"),
	"name",
	String,
)

// Literal

func MakeEnumValue(id ID, name string) EnumValueL {
	return EnumValueL{
		id,
		name,
	}
}

type EnumValueL struct {
	ID_   ID     `record:"id"`
	Name_ string `record:"name"`
}

func (lit EnumValueL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, EnumValueModel)
	return []Record{rec}, []Link{}
}

func (lit EnumValueL) ID() ID {
	return lit.ID_
}

func (lit EnumValueL) Name() string {
	return lit.Name_
}

// Dynamic

type enumValue struct {
	rec Record
}

func (ev *enumValue) ID() ID {
	return ev.rec.ID()
}

func (ev *enumValue) Name() string {
	return evName.MustGet(ev.rec).(string)
}
