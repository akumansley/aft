package data

type FieldType int

const (
	Int FieldType = iota
	String
	Text
	Float
	Enum
)

type Field struct {
	Name string    `json:"name"`
	Type FieldType `json:"type"`
}

type Object struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}
