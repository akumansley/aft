package db

var ComputedAttributeModel = ModelL{
	ID:         MakeID("29d88992-1855-4abc-968c-cf06e0979420"),
	Name:       "computedAttribute",
	Attributes: []AttributeL{},
}

var caName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("dc8ee712-f8a1-4b72-bbf7-f17d74beb796"),
	Datatype: String,
}

var ComputedAttributeGetter = RelationshipL{
	Name:   "getter",
	ID:     MakeID("2eaec801-07df-4d9f-a7b0-e1ab2a72004d"),
	Source: ComputedAttributeModel,
	Target: NativeFunctionModel,
	Multi:  false,
}

type ComputedAttributeLoader struct{}

func (l ComputedAttributeLoader) ProvideModel() ModelL {
	return ComputedAttributeModel
}

func (l ComputedAttributeLoader) Load(tx Tx, rec Record) Attribute {
	return &computedAttr{rec, tx}
}

type ComputedAttributeL struct {
	ID     ID     `record:"id"`
	Name   string `record:"name"`
	Getter NativeFunctionL
}

func (lit ComputedAttributeL) GetID() ID {
	return lit.ID
}

func (lit ComputedAttributeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, ComputedAttributeModel)
	dtl := Link{rec.ID(), lit.Getter.GetID(), ComputedAttributeGetter}
	return []Record{rec}, []Link{dtl}
}
