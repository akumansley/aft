package db

// Model

var EnumModel = ModelL{
	ID:   MakeID(""),
	Name: "enum",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID(""),
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

type EnumL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Values []EnumValueL
}

func (lit EnumL) GetID() ID {
	return lit.ID
}

func (lit EnumL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, EnumModel)

	recs = append(recs, rec)
	for _, a := range lit.Values {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)

		links = append(links, Link{rec.ID(), a.GetID(), EnumEnumValues})
	}

	return
}

func (lit EnumL) AsDatatype() Datatype {
	return eBox{lit}
}

// "Boxed" literal

type eBox struct {
	EnumL
}

func (e eBox) ID() ID {
	return e.EnumL.ID
}
func (e eBox) Name() string {
	return e.EnumL.Name
}
func (e eBox) Storage() EnumValue {
	return UUIDStorage.AsEnumValue()
}

func (e eBox) FromJSON() (Function, error) {
	panic("not implemented")
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
