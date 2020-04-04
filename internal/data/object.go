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

type Atom struct {
	Attributes       map[string]interface{} `json:"attributes"`
	RelationshipData map[string]interface{} `json:"relationships"`
}

func FromJson(jsonObject string) []Atom {

}

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
	// maybe we don't use the json tags
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

var ObjectObject Object = Object{
	Id:   "", // uuid given at db instantiation
	Name: "Object",
	Relationships: []Relationship{
		Relationship{
			Id:          "",
			Name:        "relationships",
			Cardinality: Many,
			Target:      "relationships",
		},
		Relationship{
			Id:          "",
			Name:        "attributes",
			Cardinality: Many,
			Target:      "attributes",
		},
	},
	IsSystemObject: true,
}

var ObjectIdAttribute Attribute = Attribute{
	Id:   "",
	Name: "id",
	Type: UUID,
}

var ObjectNameAttribute Attribute = Attribute{
	Id:   "",
	Name: "name",
	Type: String,
}

var ObjectRelationshipRelationship Relationship = Relationship{
	Id:          "",
	Name:        "relationships",
	Cardinality: Many,
	Target:      "relationships",
}

var ObjectAttributeRelationship Relationship = Relationship{
	Id:          "",
	Name:        "attributes",
	Cardinality: Many,
	Target:      "attributes",
}

var AttributeObject Object = Object{
	Id:             "", // uuid given at db instantiation
	Name:           "Attribute",
	IsSystemObject: true,
}

var AttributeIdAttribute Attribute = Attribute{
	Name: "id",
	Type: UUID,
}
var AttributeNameAttribute Attribute = Attribute{
	Name: "name",
	Type: String,
}
var AttributeTypeAttribute Attribute = Attribute{
	Name: "type",
	Type: Int,
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
