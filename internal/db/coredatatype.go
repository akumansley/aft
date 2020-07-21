package db

import (
	"github.com/google/uuid"
)

// Model

var CoreDatatypeModel = MakeModel(
	MakeID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	"coredatatype",
	[]AttributeL{
		cdStoredAs,
		cdName,
	},
	[]RelationshipL{
		DatatypeValidator,
	},
	[]ConcreteInterfaceL{},
)

var cdStoredAs = MakeConcreteAttribute(
	MakeID("523edf8d-6ea5-4745-8182-98165a75d4da"),
	"storedAs",
	StoredAs,
)

var cdName = MakeConcreteAttribute(
	MakeID("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
	"name",
	String,
)

var DatatypeValidator = MakeConcreteRelationship(
	MakeID("353a1d40-d292-47f6-b45c-06b059bed882"),
	"validator",
	false,
	FunctionInterface,
)

// Loader

type CoreDatatypeLoader struct{}

func (l CoreDatatypeLoader) ProvideModel() ModelL {
	return CoreDatatypeModel
}

func (l CoreDatatypeLoader) Load(tx Tx, rec Record) Datatype {
	return &coreDatatype{rec, tx}
}

// Literal

func MakeCoreDatatype(id ID, name string, storedAs EnumValueL, validator FunctionL) CoreDatatypeL {
	return CoreDatatypeL{
		id,
		name,
		storedAs,
		validator,
	}
}

type CoreDatatypeL struct {
	ID_        ID         `record:"id"`
	Name_      string     `record:"name"`
	StoredAs_  EnumValueL `record:"storedAs"`
	Validator_ FunctionL
}

func (lit CoreDatatypeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, CoreDatatypeModel)
	dtl := Link{rec.ID(), lit.Validator_.ID(), DatatypeValidator}
	return []Record{rec}, []Link{dtl}
}

func (lit CoreDatatypeL) ID() ID {
	return lit.ID_
}
func (lit CoreDatatypeL) Name() string {
	return lit.Name_
}

func (lit CoreDatatypeL) Storage() EnumValue {
	return lit.StoredAs_

}

func (lit CoreDatatypeL) FromJSON() (Function, error) {
	return lit.Validator_, nil
}

// Dynamic

type coreDatatype struct {
	rec Record
	tx  Tx
}

func (cd *coreDatatype) ID() ID {
	return cd.rec.ID()
}

func (cd *coreDatatype) Name() string {
	return cdName.MustGet(cd.rec).(string)
}

func (cd *coreDatatype) Storage() EnumValue {
	evid := cdStoredAs.MustGet(cd.rec).(uuid.UUID)
	ev, err := cd.tx.Schema().GetEnumValueByID(ID(evid))
	if err != nil {
		panic(err)
	}
	return ev
}

func (cd *coreDatatype) FromJSON() (Function, error) {
	vrec, err := cd.tx.getRelatedOne(cd.rec.ID(), DatatypeValidator.ID())
	if err != nil {
		return nil, err
	}
	return cd.tx.loadFunction(vrec)
}
