package db

// Model

var ConcreteAttributeModel = ModelL{
	ID:   MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "concreteAttribute",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
	},
}

var ConcreteAttributeDatatype = ConcreteRelationshipL{
	Name:   "datatype",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: ConcreteAttributeModel,
	Target: CoreDatatypeModel,
	Multi:  false,
}

// Loader

type ConcreteAttributeLoader struct{}

func (l ConcreteAttributeLoader) ProvideModel() ModelL {
	return ConcreteAttributeModel
}

func (l ConcreteAttributeLoader) Load(tx Tx, rec Record) Attribute {
	return &concreteAttr{rec, tx}
}

// Literal

type ConcreteAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype DatatypeL
}

func (lit ConcreteAttributeL) GetID() ID {
	return lit.ID
}

func (lit ConcreteAttributeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, ConcreteAttributeModel)
	dtl := Link{rec.ID(), lit.Datatype.GetID(), ConcreteAttributeDatatype}
	return []Record{rec}, []Link{dtl}
}

func (lit ConcreteAttributeL) AsAttribute() Attribute {
	return cBox{lit}
}

// "Boxed" literal

type cBox struct {
	ConcreteAttributeL
}

func (c cBox) ID() ID {
	return c.ConcreteAttributeL.ID
}

func (c cBox) Name() string {
	return c.ConcreteAttributeL.Name
}

func (c cBox) Datatype() Datatype {
	return c.ConcreteAttributeL.Datatype.AsDatatype()
}

func (c cBox) Storage() EnumValue {
	return c.Datatype().Storage()
}

func (c cBox) Getter() Function {
	panic("Not implemented")
}

func (c cBox) Setter() Function {
	panic("Not implemented")
}

func (c cBox) Get(Record) (interface{}, error) {
	panic("Not implemented")
}

func (c cBox) MustGet(Record) interface{} {
	panic("Not implemented")
}

func (c cBox) Set(interface{}, Record) error {
	panic("Not implemented")
}

// Dynamic

type concreteAttr struct {
	rec Record
	tx  Tx
}

func (a *concreteAttr) ID() ID {
	return a.rec.ID()
}

func (a *concreteAttr) Name() string {
	model := a.tx.Schema().GetModelByID(a.rec.Model().ID())
	nameAttr, err := model.AttributeByName("name")
	if err != nil {
		panic(err)
	}
	return nameAttr.MustGet(a.rec).(string)
}

func (a *concreteAttr) Datatype() Datatype {
	ad, _ := a.tx.Schema().GetRelationshipByID(ConcreteAttributeDatatype.ID)
	dt, _ := a.tx.GetRelatedOne(a.ID(), ad)
	return &coreDatatype{dt, a.tx}
}

func (a *concreteAttr) Storage() EnumValue {
	return a.Datatype().Storage()
}

func (a *concreteAttr) Get(rec Record) (interface{}, error) {
	return rec.get(a.Name())
}

func (a *concreteAttr) MustGet(rec Record) interface{} {
	v, err := a.Get(rec)
	if err != nil {
		panic(err)
	}
	return v
}

func (a *concreteAttr) Set(v interface{}, rec Record) error {
	f, err := a.Datatype().FromJSON()
	if err != nil {
		return err
	}
	parsed, err := f.Call(v)
	if err != nil {
		return err
	}
	rec.set(a.Name(), parsed)
	return err
}
