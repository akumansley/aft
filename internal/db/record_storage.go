package db

var StoredAs = MakeEnum(
	MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	"storedAs",
	[]EnumValueL{
		BoolStorage,
		IntStorage,
		StringStorage,
		BytesStorage,
		FloatStorage,
		UUIDStorage,
		NotStored,
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

var BytesStorage = MakeEnumValue(
	MakeID("bc7a618f-e87a-4044-a451-9e239212fe2e"),
	"bytes",
)

var FloatStorage = MakeEnumValue(
	MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
	"float",
)

var UUIDStorage = MakeEnumValue(
	MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
	"uuid",
)
