package db

import (
	"fmt"
)

// Model

var ConcreteAttributeModel = MakeModel(
	MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	"concreteAttribute",
	[]AttributeL{
		caName,
	},
	[]RelationshipL{
		ConcreteAttributeDatatype,
	},
	[]ConcreteInterfaceL{},
)

var caName = MakeConcreteAttribute(
	MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
	"name",
	String,
)

var ConcreteAttributeDatatype = MakeConcreteRelationship(
	MakeID("b503d842-7dff-48d8-90dd-398d7f9e9db3"),
	"datatype",
	false,
	DatatypeInterface,
)

var ConcreteAttributeOfModel = MakeReverseRelationship(
	MakeID("364c5d26-f0b8-4b86-ab4c-35bc5c2ac00e"),
	"model",
	ModelAttributes,
)

// breaks init loop
func init() {
	ConcreteAttributeModel.Relationships_ = append(ConcreteAttributeModel.Relationships_, ConcreteAttributeOfModel)
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

func MakeConcreteAttribute(id ID, name string, datatype DatatypeL) ConcreteAttributeL {
	return ConcreteAttributeL{
		id, name, datatype,
	}
}

type ConcreteAttributeL struct {
	ID_       ID     `record:"id"`
	Name_     string `record:"name"`
	Datatype_ DatatypeL
}

func (lit ConcreteAttributeL) MarshalDB(b *Builder) ([]Record, []Link) {
	rec := MarshalRecord(b, lit)
	dtl := Link{rec.ID(), lit.Datatype_.ID(), ConcreteAttributeDatatype}
	return []Record{rec}, []Link{dtl}
}

func (lit ConcreteAttributeL) ID() ID {
	return lit.ID_
}

func (lit ConcreteAttributeL) InterfaceID() ID {
	return ConcreteAttributeModel.ID()
}

func (lit ConcreteAttributeL) Load(tx Tx) Attribute {
	attr, err := tx.Schema().GetAttributeByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return attr
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
	return a.rec.MustGet("name").(string)
}

func (a *concreteAttr) Datatype() Datatype {
	dtrec, err := a.tx.getRelatedOne(a.ID(), ConcreteAttributeDatatype.ID())
	if err != nil {
		err = fmt.Errorf("%w: %v.Datatype", err, a.rec)
		panic(err)
	}

	dt, err := a.tx.Schema().loadDatatype(dtrec)
	if err != nil {
		panic(err)
	}
	return dt
}

func (a *concreteAttr) Storage() EnumValue {
	return a.Datatype().Storage()
}

func (a *concreteAttr) Get(rec Record) (interface{}, error) {
	return rec.Get(a.Name())
}

func (a *concreteAttr) MustGet(rec Record) interface{} {
	v, err := a.Get(rec)
	if err != nil {
		panic(err)
	}
	return v
}

func (a *concreteAttr) Set(rec Record, v interface{}) error {
	f, err := a.Datatype().FromJSON()
	if err != nil {
		fmt.Printf("HIHI1\n")
		return err
	}
	parsed, err := f.Call([]interface{}{v, rec})
	if err != nil {
		fmt.Printf("HIHI2\n")
		return err
	}
	rec.Set(a.Name(), parsed)
	return err
}
