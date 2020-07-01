package db

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
