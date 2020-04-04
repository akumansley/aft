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
	Id   string    `json:"id"`
	Name string    `json:"name"`
	Type FieldType `json:"type"`
}

func (a Attribute) GetId() string {
	return a.Id
}

type Relationship struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Cardinality Cardinality `json:"cardinality"`
	Target      string      `json:"target"`
	// maybe in the future, other application-level concepts
	// like "ownership" and polymorphism
}

func (r Relationship) GetId() string {
	return r.Id
}

type Object struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	Attributes     []Attribute    `json:"attributes"`
	Relationships  []Relationship `json:"relationships"`
	IsSystemObject bool           `json:-`
}

func (o Object) GetId() string {
	return o.Id
}

// some bootstrapping here..
var ObjectObject Object = Object{
	Id:   "", // uuid given at db instantiation
	Name: "Object",
	Attributes: []Attribute{
		Attribute{
			Name: "id",
			Type: UUID,
		},
		Attribute{
			Name: "name",
			Type: String,
		},
	},
	Relationships:  []Relationship{},
	IsSystemObject: true,
}

var AttributeObject Object = Object{
	Id:   "", // uuid given at db instantiation
	Name: "Attribute",
	Attributes: []Attribute{
		Attribute{
			Name: "id",
			Type: UUID,
		},
		Attribute{
			Name: "name",
			Type: String,
		},
		Attribute{
			Name: "type",
			Type: Int,
		},
	},
	IsSystemObject: true,
}

var RelationshipObject Object = Object{
	Id:   "", // uuid given at db instantiation
	Name: "Relationship",
	Attributes: []Attribute{
		Attribute{
			Name: "id",
			Type: UUID,
		},
		Attribute{
			Name: "name",
			Type: String,
		},
		Attribute{
			Name: "cardinality",
			Type: Int,
		},
		Attribute{
			Name: "target",
			Type: String,
		},
	},
	IsSystemObject: true,
}
