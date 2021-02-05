package db

import "context"

type Model interface {
	ID() ID
	Name() string
	Implements(Tx) ([]Interface, error)
	Attributes(Tx) ([]Attribute, error)
	AttributeByName(Tx, string) (Attribute, error)
	Relationships(Tx) ([]Relationship, error)
	Targeted(Tx) ([]Relationship, error)
	RelationshipByName(Tx, string) (Relationship, error)
}

type Interface interface {
	ID() ID
	Name() string
	Attributes(Tx) ([]Attribute, error)
	AttributeByName(Tx, string) (Attribute, error)
	Relationships(Tx) ([]Relationship, error)
	RelationshipByName(Tx, string) (Relationship, error)
}

type Attribute interface {
	ID() ID
	Name() string
	Get(Record) (interface{}, error)
	MustGet(Record) interface{}

	// set takes Tx bc of validation
	Set(Tx, Record, interface{}) error
	Storage(Tx) EnumValue
	Datatype(Tx) Datatype
}

type Relationship interface {
	ID() ID
	Name() string
	Multi() bool

	Connect(RWTx, Record, Record) error
	Disconnect(RWTx, Record, Record) error

	LoadOne(Tx, Record) (Record, error)
	LoadMany(Tx, Record) ([]Record, error)
	LoadManyReverse(Tx, Record) ([]Record, error)

	Source(Tx) Interface
	Target(Tx) Interface
}

type Datatype interface {
	ID() ID
	Name() string
	Storage(Tx) EnumValue
	FromJSON(Tx) (Function, error)
}

type Function interface {
	ID() ID
	Name() string
	Arity() int
	FuncType(Tx) EnumValue
	Call(context.Context, []interface{}) (interface{}, error)
}

type FunctionLoader interface {
	ProvideModel() ModelL
	Load(Record) Function
}

type AttributeLoader interface {
	ProvideModel() ModelL
	Load(Record) Attribute
}

type RelationshipLoader interface {
	ProvideModel() ModelL
	Load(Record) Relationship
}

type DatatypeLoader interface {
	ProvideModel() ModelL
	Load(Record) Datatype
}

type InterfaceLoader interface {
	ProvideModel() ModelL
	Load(Record) Interface
}

type EnumValue interface {
	ID() ID
	Name() string
}
