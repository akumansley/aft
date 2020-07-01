package db

type RelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Multi  bool   `record:"multi"`
	Target Literal
	Source Literal
}

func (lit RelationshipL) GetID() ID {
	return lit.ID
}

func (lit RelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, RelationshipModel)
	recs = append(recs, rec)
	source := Link{rec.ID(), lit.Source.GetID(), RelationshipSource}
	target := Link{rec.ID(), lit.Target.GetID(), RelationshipTarget}
	links = []Link{source, target}
	return
}
