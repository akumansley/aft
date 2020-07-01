package db

type rBox struct {
	RelationshipL
}

type RelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Multi  bool   `record:"multi"`
	Target Model
	Source Model
}

func (lit RelationshipL) AsRelationship() Relationship {
	return rBox{lit}
}

func (r rBox) ID() ID {
	return r.RelationshipL.ID
}

func (r rBox) Name() string {
	return r.RelationshipL.Name
}

func (r rBox) Multi() bool {
	return r.RelationshipL.Multi
}

func (r rBox) Source() Model {
	return r.RelationshipL.Source
}

func (r rBox) Target() Model {
	return r.RelationshipL.Target
}
