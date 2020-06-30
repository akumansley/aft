package db

type Model interface {
	ID() ModelID
	Name() string
	Attributes() ([]Attribute, error)
	Relationships() ([]Relationship, error)
}

type Interface interface {
	ID() ID
	Name() string
	Attributes() ([]Attribute, error)
}

type Attribute interface {
	ID() ID
	Name() string
	Datatype() Datatype
}

type Relationship interface {
	ID() ID
	Name() string
	Multi() bool
	Source() Model
	Target() Model
}

type Datatype interface {
	ID() ID
	Name() string
	Storage() StorageEnumValue
	FromJSON() (Function, error)
}

type Function interface {
	ID() ID
	Name() string
	Runtime() RuntimeEnumValue
	FunctionSignature() FunctionSignatureEnumValue
	Call(interface{}) (interface{}, error)
}
