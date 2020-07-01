package db

import (
	"awans.org/aft/internal/datatypes"
)

var ModelModel = ModelL{
	ID:   MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var AttributeModel = ModelL{
	ID:   MakeModelID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		}.AsAttribute(),
		AttributeL{ //todo remove hack
			Name:     "datatypeID",
			ID:       MakeID("bfeefcbf-b9f7-44e6-9951-134755f7e1cd"),
			Datatype: UUID,
		}.AsAttribute(),
	},
}.AsModel()

var RelationshipModel = ModelL{
	ID:   MakeModelID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		}.AsAttribute(),
		AttributeL{
			Name:     "multi",
			ID:       MakeID("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Bool,
		}.AsAttribute(),
	},
}.AsModel()

var ModelAttributes = RelationshipL{
	Name:   "attributes",
	ID:     MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	Source: ModelModel,
	Target: AttributeModel,
	Multi:  true,
}.AsRelationship()

var AttributeDatatype = RelationshipL{
	Name:   "datatype",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: AttributeModel,
	Target: DatatypeModel,
	Multi:  false,
}.AsRelationship()

var RelationshipSource = RelationshipL{
	Name:   "source",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: RelationshipModel,
	Target: ModelModel, // model
	Multi:  false,
}.AsRelationship()

var RelationshipTarget = RelationshipL{
	Name:   "target",
	ID:     MakeID("e194f9bf-ea7a-4c78-a179-bdf9c044ac3c"),
	Source: RelationshipModel,
	Target: ModelModel,
	Multi:  false,
}.AsRelationship()

var DatatypeModel = ModelL{
	ID:   MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name: "datatype",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		}.AsAttribute(),
		AttributeL{
			Name:     "storedAs",
			ID:       MakeID("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: StoredAs,
		}.AsAttribute(),
		AttributeL{
			Name:     "enum",
			ID:       MakeID("931050f5-0022-4be2-87fb-d69537877a87"),
			Datatype: Bool,
		}.AsAttribute(),
		AttributeL{
			Name:     "native",
			ID:       MakeID("db56571e-1939-45f1-b122-9ecb8ad9fd7e"),
			Datatype: Bool,
		}.AsAttribute(),
	},
}.AsModel()

var NativeFunctionModel = ModelL{
	ID:   MakeModelID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
			Datatype: String,
		}.AsAttribute(),
		AttributeL{
			Name:     "runtime",
			ID:       MakeID("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Runtime,
		}.AsAttribute(),
		AttributeL{
			Name:     "functionSignature",
			ID:       MakeID("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
			Datatype: FunctionSignature,
		}.AsAttribute(),
		AttributeL{
			Name:     "code",
			ID:       MakeID("80b3055b-08ad-41fe-b562-4a493bb6db36"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var EnumValueModel = ModelL{
	ID:   MakeModelID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"),
	Name: "enumValue",
	Attributes: []Attribute{
		AttributeL{
			Name:     "name",
			ID:       MakeID("5803e350-48f8-448d-9901-7c80f45c775b"),
			Datatype: String,
		}.AsAttribute(),
		AttributeL{
			Name:     "value",
			ID:       MakeID("9dabda3c-57af-4814-909d-8c2299c236e8"),
			Datatype: Int,
		}.AsAttribute(),
	},
}.AsModel()

var DatatypeValidator = RelationshipL{
	ID:     MakeID("353a1d40-d292-47f6-b45c-06b059bed882"),
	Name:   "validator",
	Source: DatatypeModel,       // datatype
	Target: NativeFunctionModel, // code
	Multi:  false,
}.AsRelationship()

var DatatypeEnumValues = RelationshipL{
	ID:     MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
	Source: DatatypeModel, // datatype
	Name:   "enumValues",
	Multi:  true,
	Target: EnumValueModel,
}

var boolValidator = NativeFunctionL{
	ID:                MakeID("8e806967-c462-47af-8756-48674537a909"),
	Name:              "bool",
	Function:          datatypes.BoolFromJSON,
	FunctionSignature: FromJSON,
}.AsFunction()

var intValidator = NativeFunctionL{
	Name:              "int",
	ID:                MakeID("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Function:          datatypes.IntFromJSON,
	FunctionSignature: FromJSON,
}.AsFunction()

// var enumValidator = NativeFunctionL{
// 	Name:              "enum",
// 	ID:                MakeID("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
// 	Runtime:           Native,
// 	Function:          datatypes.EnumFromJSON,
// 	FunctionSignature: FromJSON,
// }.AsFunction()

var stringValidator = NativeFunctionL{
	Name:              "string",
	ID:                MakeID("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Function:          datatypes.StringFromJSON,
	FunctionSignature: FromJSON,
}.AsFunction()

// var textValidator = NativeFunctionL{
// 	Name:              "text",
// 	ID:                MakeID("9f10ac9f-afd2-423a-8857-d900a0c97563"),
// 	Runtime:           Native,
// 	Function:          datatypes.TextFromJSON,
// 	FunctionSignature: FromJSON,
// }.AsFunction()

var uuidValidator = NativeFunctionL{
	Name:              "uuid",
	ID:                MakeID("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Function:          datatypes.UUIDFromJSON,
	FunctionSignature: FromJSON,
}.AsFunction()

var floatValidator = NativeFunctionL{
	Name:              "float",
	ID:                MakeID("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Function:          datatypes.FloatFromJSON,
	FunctionSignature: FromJSON,
}.AsFunction()

var Bool = CoreDatatypeL{
	ID:        MakeID("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:      "bool",
	Validator: boolValidator,
	StoredAs:  BoolStorage,
}.AsDatatype()

var Int = CoreDatatypeL{
	ID:        MakeID("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:      "int",
	Validator: intValidator,
	StoredAs:  IntStorage,
}.AsDatatype()

var String = CoreDatatypeL{
	ID:        MakeID("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:      "string",
	Validator: stringValidator,
	StoredAs:  StringStorage,
}.AsDatatype()

var UUID = CoreDatatypeL{
	ID:        MakeID("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:      "uuid",
	Validator: uuidValidator,
	StoredAs:  UUIDStorage,
}.AsDatatype()

var Float = CoreDatatypeL{
	ID:        MakeID("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:      "float",
	Validator: floatValidator,
	StoredAs:  FloatStorage,
}.AsDatatype()

var Runtime = Enum{
	ID:   MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name: "runtime",
}

var FunctionSignature = Enum{
	ID:   MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"),
	Name: "functionSignature",
}

var StoredAs = Enum{
	ID:   MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	Name: "storedAs",
}

var Native = RuntimeEnumValue{
	EnumValue{
		ID:   MakeID("cecf8eac-d3be-4ca0-927a-127763d465b1"),
		Name: "native",
		// Datatype: MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"), //Runtime ID
	},
}

var Starlark = RuntimeEnumValue{
	EnumValue{
		ID:   MakeID("c0036590-8227-46cb-8cf9-689dd17616a3"),
		Name: "starlark",
		// Datatype: MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"), //Runtime ID
	},
}

var FromJSON = FunctionSignatureEnumValue{
	EnumValue{
		ID:   MakeID("508ba2cc-ce86-4615-bc0d-fe0d085a2051"),
		Name: "fromJson",
		// Datatype: MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"), //FunctionSignature ID
	},
}

var RPC = FunctionSignatureEnumValue{
	EnumValue{
		ID:   MakeID("8decedba-555b-47ca-a232-68100fbbf756"),
		Name: "rpc",
		// Datatype: MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"), //FunctionSignature ID
	},
}

var BoolStorage = StorageEnumValue{
	EnumValue{
		ID:   MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
		Name: "bool",
		// Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var IntStorage = StorageEnumValue{
	EnumValue{
		ID:   MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
		Name: "int",
		// Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var StringStorage = StorageEnumValue{
	EnumValue{
		ID:   MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
		Name: "string",
		// Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var FloatStorage = StorageEnumValue{
	EnumValue{
		ID:   MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
		Name: "float",
		// Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var UUIDStorage = StorageEnumValue{
	EnumValue{
		ID:   MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
		Name: "uuid",
		// Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}
