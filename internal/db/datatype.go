package db

// Interface definition

var DatatypeInterface = MakeInterface(
	MakeID("3f518387-ac9b-47bd-9b66-69b9a271f0d2"),
	"datatype",
	[]AttributeL{
		dName,
	}, []InterfaceRelationshipL{DatatypeModule},
)

var dName = MakeConcreteAttribute(
	MakeID("52d1cf2c-31d8-40a7-8c5e-9f27dfae9064"),
	"name",
	String,
)

var DatatypeModule = MakeInterfaceRelationship(
	MakeID("4daa535f-0254-48ff-aa21-2e83823de0d8"),
	"module",
	false,
	ModuleModel,
)
