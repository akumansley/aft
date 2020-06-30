package db

type mlBox struct {
	Model
}

type ModelL struct {
	ID         ModelID `record:"id"`
	Name       string  `record:"name"`
	Attributes []Attribute
}

func (m mlBox) ID() ModelID {
	return m.Model.ID
}

func (m mlBox) Name() string {
	return m.Model.Name
}

func (m mlBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) Attributes() ([]Attribute, error) {
	return m.Model.Attributes, nil
}

func (s Schema) SaveModel(mL ModelL) (err error) {
	m := mlBox{mL}
	rec, err := MarshalRecord(m, mlBox{ModelModel})
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
