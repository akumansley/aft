package db

type rBox struct {
	RelationshipL
}

type RelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Multi  bool   `record:"multi"`
	Target Interface
	Source Interface
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

func (r rBox) Source() Interface {
	return r.RelationshipL.Source
}

func (r rBox) Target() Interface {
	return r.RelationshipL.Target
}
