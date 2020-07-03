package db

var ModelAttributes = ConcreteRelationshipL{
	Name:   "attributes",
	ID:     MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	Source: ModelModel,
	Target: ConcreteAttributeModel,
	Multi:  true,
}

var FunctionSignature = EnumL{
	ID:   MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"),
	Name: "functionSignature",
	Values: []EnumValueL{
		FromJSON,
		RPC,
	},
}

var StoredAs = EnumL{
	ID:   MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	Name: "storedAs",
	Values: []EnumValueL{
		BoolStorage,
		IntStorage,
		StringStorage,
		FloatStorage,
		UUIDStorage,
	},
}

var FromJSON = EnumValueL{
	ID:   MakeID("508ba2cc-ce86-4615-bc0d-fe0d085a2051"),
	Name: "fromJson",
}

var RPC = EnumValueL{
	ID:   MakeID("8decedba-555b-47ca-a232-68100fbbf756"),
	Name: "rpc",
}

var NotStored = EnumValueL{
	ID:   MakeID("e0f86fe9-10ea-430b-a393-b01957a3eabf"),
	Name: "notStored",
}

var BoolStorage = EnumValueL{
	ID:   MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
	Name: "bool",
}

var IntStorage = EnumValueL{
	ID:   MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
	Name: "int",
}

var StringStorage = EnumValueL{
	ID:   MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
	Name: "string",
}

var FloatStorage = EnumValueL{
	ID:   MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
	Name: "float",
}

var UUIDStorage = EnumValueL{
	ID:   MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
	Name: "uuid",
}
