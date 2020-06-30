package db

type RelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Multi  bool   `record:"multi"`
	Target ModelL
	Source ModelL
}

func (s Schema) SaveRelationship(r Relationship) (err error) {
	rec, err := MarshalRecord(r, RelationshipModel)
	if err != nil {
		return
	}
	s.tx.Insert(rec)
	s.tx.Connect(rec.ID(), ID(r.Source().ID()), RelationshipSource)
	s.tx.Connect(rec.ID(), ID(r.Target().ID()), RelationshipTarget)
	return
}
