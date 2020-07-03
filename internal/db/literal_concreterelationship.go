package db

var ConcreteRelationshipModel = ModelL{
	ID:   MakeID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		},
		ConcreteAttributeL{
			Name:     "multi",
			ID:       MakeID("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Bool,
		},
	},
}

var ConcreteRelationshipSource = ConcreteRelationshipL{
	Name:   "source",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: ConcreteRelationshipModel,
	Target: ModelModel,
	Multi:  false,
}

var ConcreteRelationshipTarget = ConcreteRelationshipL{
	Name:   "target",
	ID:     MakeID("e194f9bf-ea7a-4c78-a179-bdf9c044ac3c"),
	Source: ConcreteRelationshipModel,
	Target: ModelModel,
	Multi:  false,
}

type ConcreteRelationshipL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Multi  bool   `record:"multi"`
	Target Literal
	Source Literal
}

func (lit ConcreteRelationshipL) GetID() ID {
	return lit.ID
}

func (lit ConcreteRelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ConcreteRelationshipModel)
	recs = append(recs, rec)
	source := Link{rec.ID(), lit.Source.GetID(), ConcreteRelationshipSource}
	target := Link{rec.ID(), lit.Target.GetID(), ConcreteRelationshipTarget}
	links = []Link{source, target}
	return
}
