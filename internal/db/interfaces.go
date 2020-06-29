package db

type Model interface {
	ID() ID
	Name() string
	Implements() ([]Interface, error)
	Attributes() ([]Attribute, error)
	AttributeByName(string) (Attribute, error)
	Relationships() ([]Relationship, error)
	RelationshipByName(string) (Relationship, error)
}

type Interface interface {
	ID() ID
	Name() string
	Attributes() ([]Attribute, error)
	AttributeByName(string) (Attribute, error)
	Relationships() ([]Relationship, error)
	RelationshipByName(string) (Relationship, error)
}

type Attribute interface {
	ID() ID
	Name() string
	Storage() EnumValue
	Datatype() Datatype
	Get(Record) (interface{}, error)
	MustGet(Record) interface{}
	Set(Record, interface{}) error
}

type Relationship interface {
	ID() ID
	Name() string
	Multi() bool
	LoadOne(Record) (Record, error)
	LoadMany(Record) ([]Record, error)
	Source() Interface
	Target() Interface
}

type Datatype interface {
	ID() ID
	Name() string
	Storage() EnumValue
	FromJSON() (Function, error)
}

type Function interface {
	ID() ID
	Name() string
	FunctionSignature() EnumValue
	Call(interface{}) (interface{}, error)
}

type FunctionLoader interface {
	ProvideModel() ModelL
	Load(Tx, Record) Function
}

type AttributeLoader interface {
	ProvideModel() ModelL
	Load(Tx, Record) Attribute
}

type RelationshipLoader interface {
	ProvideModel() ModelL
	Load(Tx, Record) Relationship
}

type DatatypeLoader interface {
	ProvideModel() ModelL
	Load(Tx, Record) Datatype
}

type EnumValue interface {
	ID() ID
	Name() string
}
