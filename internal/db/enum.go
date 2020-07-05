package db

// Model

var EnumModel = ModelL{
	ID:   MakeID("7f0d30ae-83c7-4d89-a134-e8ac326321e6"),
	Name: "enum",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("57ff4302-41b4-4866-8a5c-bc1c264ffba4"),
			Datatype: String,
		},
	},
}

var EnumEnumValues = ConcreteRelationshipL{
	ID:     MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
	Name:   "enumValues",
	Source: EnumModel,
	Target: EnumValueModel,
	Multi:  true,
}

// Loader

type EnumDatatypeLoader struct{}

func (l EnumDatatypeLoader) ProvideModel() ModelL {
	return EnumModel
}

func (l EnumDatatypeLoader) Load(tx Tx, rec Record) Datatype {
	return &enum{rec, tx}
}

// Literal

func MakeEnum(id ID, name string, values []EnumValueL) EnumL {
	return EnumL{
		id,
		name,
		values,
	}
}

type EnumL struct {
	ID_     ID     `record:"id"`
	Name_   string `record:"name"`
	Values_ []EnumValueL
}

func (lit EnumL) GetID() ID {
	return lit.ID_
}

func (lit EnumL) ID() ID {
	return lit.ID_
}

func (lit EnumL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, EnumModel)

	recs = append(recs, rec)
	for _, a := range lit.Values_ {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)

		links = append(links, Link{rec.ID(), a.GetID(), EnumEnumValues})
	}

	return
}

func (e EnumL) Name() string {
	return e.Name_
}
func (e EnumL) Storage() EnumValue {
	return UUIDStorage
}

func (e EnumL) FromJSON() (Function, error) {
	// TODO write a proper enumvalidator
	return uuidValidator, nil
}

// Dynamic

type enum struct {
	rec Record
	tx  Tx
}

func (cd *enum) ID() ID {
	return cd.rec.ID()
}

func (cd *enum) Name() string {
	panic("Not Implemented")
	// return cdName.AsAttribute().MustGet(cd.rec).(string)
}

func (cd *enum) Storage() EnumValue {
	panic("Not Implemented")
	// evid := cdStoredAs.AsAttribute().MustGet(cd.rec).(uuid.UUID)
	// ev, err := cd.tx.Schema().GetEnumValueByID(ID(evid))
	// if err != nil {
	// 	panic(err)
	// }
	// return ev
}

func (cd *enum) FromJSON() (Function, error) {
	// vrec, _ := cd.tx.GetRelatedOne(cd.rec.ID(), DatatypeValidator)
	panic("Not Implemented")
}
