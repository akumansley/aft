package data

type FieldType int

const (
	Int FieldType = iota
	String
	Text
	Float
	Enum
	UUID
)

type Cardinality int

const (
	One Cardinality = iota
	Many
)

type Attribute struct {
	Type FieldType
}

type Relationship struct {
	Cardinality Cardinality
	Target      string
}

type Model struct {
	Name          string
	Attributes    map[string]Attribute
	Relationships map[string]Relationship
}
