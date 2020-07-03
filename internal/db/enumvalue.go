package db

// Model

var EnumValueModel = ModelL{
	ID:         MakeID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"),
	Name:       "enumValue",
	Attributes: []AttributeL{},
}

var evName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("5803e350-48f8-448d-9901-7c80f45c775b"),
	Datatype: String,
}

// Literal

type EnumValueL struct {
	ID   ID
	Name string
}

func (lit EnumValueL) GetID() ID {
	return lit.ID
}

func (lit EnumValueL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, EnumValueModel)
	return []Record{rec}, []Link{}
}

func (lit EnumValueL) AsEnumValue() EnumValue {
	return evBox{lit}
}

// "Boxed" literal

type evBox struct {
	EnumValueL
}

func (e evBox) ID() ID {
	return e.EnumValueL.ID
}

func (e evBox) Name() string {
	return e.EnumValueL.Name
}

// Dynamic

type enumValue struct {
	rec Record
}

func (ev *enumValue) ID() ID {
	return ev.rec.ID()
}

func (ev *enumValue) Name() string {
	return evName.AsAttribute().MustGet(ev.rec).(string)
}
