package db

type Model interface {
	ID() ID
	Name() string
	Attributes() ([]Attribute, error)
	AttributeByName(string) (Attribute, error)
	Relationships() ([]Relationship, error)
	RelationshipByName(string) (Relationship, error)
	Interfaces() ([]Interface, error)
}

type Interface interface {
	ID() ID
	Name() string
	Attributes() ([]Attribute, error)
	Relationships() ([]Relationship, error)
}

type Attribute interface {
	ID() ID
	Name() string
	Storage() EnumValue
	Get(Record) (interface{}, error)
	MustGet(Record) interface{}
	Set(interface{}, Record) error
}

type Relationship interface {
	ID() ID
	Name() string
	Multi() bool
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

// type Enum interface {
// 	Values() []EnumValue
// }

type EnumValue interface {
	ID() ID
	Name() string
}
