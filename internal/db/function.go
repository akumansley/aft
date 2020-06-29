package db

var FunctionInterface = MakeInterface(
	MakeID("6f55b11e-be7f-4f34-a6ac-1e42d1cd943e"),
	"function",
	[]AttributeL{
		fName,
		fFuncSig,
	}, []RelationshipL{},
)

var fFuncSig = MakeConcreteAttribute(
	MakeID("333ab360-c652-4824-9518-d9aaaaa6d3be"),
	"functionSignature",
	FunctionSignature,
)

var fName = MakeConcreteAttribute(
	MakeID("048d6151-d80f-44ab-9c77-9ebe70af5b74"),
	"name",
	String,
)

var FunctionSignature = MakeEnum(
	MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"),
	"functionSignature",
	[]EnumValueL{
		FromJSON,
		RPC,
	},
)
var FromJSON = MakeEnumValue(
	MakeID("508ba2cc-ce86-4615-bc0d-fe0d085a2051"),
	"fromJson",
)

var RPC = MakeEnumValue(
	MakeID("8decedba-555b-47ca-a232-68100fbbf756"),
	"rpc",
)
