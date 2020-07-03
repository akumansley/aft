package db

type mlBox struct {
	ModelL
}

type ModelL struct {
	ID         ID     `record:"id"`
	Name       string `record:"name"`
	Attributes []AttributeL
}

func (lit ModelL) GetID() ID {
	return lit.ID
}

func (lit ModelL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ModelModel)
	recs = append(recs, rec)
	for _, a := range lit.Attributes {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{rec.ID(), a.GetID(), ModelAttributes})
	}
	return
}

func (lit ModelL) AsModel() Model {
	return mlBox{lit}
}

func (m mlBox) ID() ID {
	return m.ModelL.ID
}

func (m mlBox) Lit() ModelL {
	return m.ModelL
}

func (m mlBox) Name() string {
	return m.ModelL.Name
}

func (m mlBox) Interfaces() ([]Interface, error) {
	panic("Not implemented")
}

func (m mlBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) RelationshipByName(name string) (Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) Attributes() ([]Attribute, error) {
	var attrs []Attribute
	for _, a := range m.ModelL.Attributes {
		attrs = append(attrs, a.AsAttribute())
	}
	return attrs, nil
}

func (m mlBox) AttributeByName(name string) (a Attribute, err error) {
	panic("Not implemented")
}
