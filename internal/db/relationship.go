package db

// Interface definition

var RelationshipInterface = MakeInterface(
	MakeID("4638099a-eeb1-4ee3-842d-27f14ad662f2"),
	"relationship",
	[]AttributeL{
		rName,
	}, []RelationshipL{},
)

var rName = MakeConcreteAttribute(
	MakeID("6e23071a-bca3-4b7a-9141-f39a28a057b7"),
	"name",
	String,
)

// var rMulti = MakeConcreteAttribute(
// 	MakeID("370aa3b8-004f-4291-b934-7b752dceb330"),
// 	"multi",
// 	Bool,
// )

// var ConcreteRelationshipSource = MakeReverseRelationship(
// 	MakeID("a84f4737-61b0-44d2-85f9-5250b8110b62"),
// 	"source",
// 	ModelRelationships,
// )
