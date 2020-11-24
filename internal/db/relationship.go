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
