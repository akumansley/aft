package db

type mlBox struct {
	Model
}

type ModelL struct {
	ID         ModelID `record:"id"`
	Name       string  `record:"name"`
	Attributes []Attribute
}

func (lit ModelL) AsModel() Model {
	return mlBox{lit}
}

func (m mlBox) ID() ModelID {
	return m.ModelL.ID
}

func (m mlBox) Name() string {
	return m.ModelL.Name
}

func (m mlBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) Attributes() ([]Attribute, error) {
	return m.ModelL.Attributes, nil
}

func (s Schema) SaveModel(m Model) (err error) {
	rec, err := MarshalRecord(m, ModelModel)
	if err != nil {
		return
	}

	s.tx.Insert(rec)
	for _, a := range m.Attributes() {
		var ar Record
		ar, err = MarshalRecord(a, AttributeModel)
		if err != nil {
			return
		}
		s.tx.Insert(ar)
		s.tx.Connect(rec.ID(), ar.ID(), ModelAttributes)
	}
	return
}
