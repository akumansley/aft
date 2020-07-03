package db

var ReverseRelationshipModel = ModelL{
	ID:   MakeID("988e0a48-3da1-40c6-8d1b-be9f44eac5dc"),
	Name: "relationship",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("0646d3e7-e8b6-4663-a18c-d9ccf25d61f6"),
			Datatype: String,
		},
	},
}

var ReverseRelationshipReferencing = ConcreteRelationshipL{
	Name:   "referencing",
	ID:     MakeID("fbca6418-da50-4737-ada1-98505dcaec6a"),
	Source: ReverseRelationshipModel,
	Target: ConcreteRelationshipModel,
	Multi:  false,
}

type ReverseRelationshipL struct {
	ID          ID     `record:"id"`
	Name        string `record:"name"`
	Referencing ConcreteRelationshipL
}

func (lit ReverseRelationshipL) GetID() ID {
	return lit.ID
}

func (lit ReverseRelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ReverseRelationshipModel)
	recs = append(recs, rec)
	return
}
