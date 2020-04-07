package model

type FieldType int

const (
	Int FieldType = iota
	String
	Text
	Float
	Enum
	UUID
)

type RelType int

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Attribute struct {
	Type FieldType
}

type Relationship struct {
	Type   RelType
	Target string
}

type Model struct {
	Name          string `boldholdIndex:"Name"`
	Attributes    map[string]Attribute
	Relationships map[string]Relationship
}
