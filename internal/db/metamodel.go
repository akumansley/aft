package db

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

var StoredAs = MakeEnum(
	MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	"storedAs",
	[]EnumValueL{
		BoolStorage,
		IntStorage,
		StringStorage,
		FloatStorage,
		UUIDStorage,
	})

var NotStored = MakeEnumValue(
	MakeID("e0f86fe9-10ea-430b-a393-b01957a3eabf"),
	"notStored",
)

var BoolStorage = MakeEnumValue(
	MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
	"bool",
)

var IntStorage = MakeEnumValue(
	MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
	"int",
)

var StringStorage = MakeEnumValue(
	MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
	"string",
)

var FloatStorage = MakeEnumValue(
	MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
	"float",
)

var UUIDStorage = MakeEnumValue(
	MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
	"uuid",
)
