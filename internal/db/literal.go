package db

type ModelL struct {
	ID_         ModelID      `record:"id"`
	Name_       string       `record:"name"`
	Attributes_ []AttributeL `rel:"attributes"`
}

func (m ModelL) ID() ModelID {
	return m.ID_
}

func (m ModelL) Name() string {
	return m.Name_
}

func (m ModelL) Relationships() ([]Relationship, error) {
	return m.Relationships_, nil
}

func (m ModelL) Attributes() ([]Attribute, error) {
	return m.Attributes_, nil
}

type AttributeL struct {
	ID   ID     `record:"id"`
	Name string `record:"name"`
}

type RelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Target ModelL `rel: "target"`
	Source ModelL `rel: "source"`
}

func MarshalRecord(v interface{}, m Model) Record {

}

func SaveModel(tx RWTx, m ModelL) {

}
