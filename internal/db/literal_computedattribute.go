package db

var ComputedAttributeModel = ModelL{
	ID:   MakeID("29d88992-1855-4abc-968c-cf06e0979420"),
	Name: "computedAttribute",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("dc8ee712-f8a1-4b72-bbf7-f17d74beb796"),
			Datatype: String,
		},
	},
}

// add a relationship to 'getter'

type ComputedAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype DatatypeL
	Getter   NativeFunctionL
}

func (lit ComputedAttributeL) GetID() ID {
	return lit.ID
}

func (lit ComputedAttributeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, ComputedAttributeModel)
	dtl := Link{rec.ID(), lit.Datatype.GetID(), AttributeDatatype}
	return []Record{rec}, []Link{dtl}
}
