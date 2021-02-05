package db

import (
	"github.com/google/uuid"
)

// Model

var CoreDatatypeModel = MakeModel(
	MakeID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	"coreDatatype",
	[]AttributeL{
		cdStoredAs,
		cdName,
	},
	[]RelationshipL{
		DatatypeValidator,
		CoreDatatypeModule,
	},
	[]ConcreteInterfaceL{DatatypeInterface},
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

var CoreDatatypeModule = MakeConcreteRelationship(
	MakeID("15a251c8-42f5-4498-9b4b-29be071a989c"),
	"module",
	false,
	ModuleModel,
)

// Loader

type CoreDatatypeLoader struct{}

func (l CoreDatatypeLoader) ProvideModel() ModelL {
	return CoreDatatypeModel
}

func (l CoreDatatypeLoader) Load(rec Record) Datatype {
	return &coreDatatype{rec}
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

func (lit CoreDatatypeL) MarshalDB(b *Builder) ([]Record, []Link) {
	rec := MarshalRecord(b, lit)
	dtl := Link{lit, lit.Validator_, DatatypeValidator}
	return []Record{rec}, []Link{dtl}
}

func (lit CoreDatatypeL) ID() ID {
	return lit.ID_
}

func (lit CoreDatatypeL) InterfaceID() ID {
	return CoreDatatypeModel.ID()
}

func (lit CoreDatatypeL) InterfaceName() string {
	return CoreDatatypeModel.Name_
}

func (lit CoreDatatypeL) Load(tx Tx) Datatype {
	dt, err := tx.Schema().GetDatatypeByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return dt
}

// Dynamic

type coreDatatype struct {
	rec Record
}

func (cd *coreDatatype) ID() ID {
	return cd.rec.ID()
}

func (cd *coreDatatype) Name() string {
	return cd.rec.MustGet("name").(string)
}

func (cd *coreDatatype) Storage(tx Tx) EnumValue {
	evid := cd.rec.MustGet("storedAs").(uuid.UUID)
	ev, err := tx.Schema().GetEnumValueByID(ID(evid))
	if err != nil {
		panic(err)
	}
	return ev
}

func (cd *coreDatatype) FromJSON(tx Tx) (Function, error) {
	vrec, err := tx.getRelatedOne(cd.rec.ID(), DatatypeValidator.ID())
	if err != nil {
		return nil, err
	}
	return tx.Schema().LoadFunction(vrec)
}
