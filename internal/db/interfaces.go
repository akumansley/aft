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
	Datatype() Datatype
	Get(Record) interface{}
	Set(interface{}, Record)
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

type Runtime interface {
	ProvideModel() Model
	Load(Tx, Record) Function
	Registered(DB)
}

// type Enum interface {
// 	Values() []EnumValue
// }

type EnumValue interface {
	ID() ID
	Name() string
}
