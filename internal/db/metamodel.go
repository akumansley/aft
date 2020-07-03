package db

import (
	"awans.org/aft/internal/datatypes"
)

var ModelModel = Model{
	ID:     MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name:   "model",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Datatype: String,
		},
		Attribute{
			Name:     "system",
			ID:       MakeID("6cd21120-7a22-4d70-9890-d5b07a63b106"),
			Datatype: Bool,
		},
	},
	LeftRelationships: []Relationship{
		ModelAttributes,
		ModelRelationshipsLeft,
		ModelRelationshipsRight,
	},
}

var AttributeModel = Model{
	ID:     MakeModelID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name:   "attribute",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
		Attribute{ //todo remove hack
			Name:     "datatypeId",
			ID:       MakeID("bfeefcbf-b9f7-44e6-9951-134755f7e1cd"),
			Datatype: UUID,
		},
	},
	RightRelationships: []Relationship{
		ModelAttributes,
	},
	LeftRelationships: []Relationship{
		AttributeDatatype,
	},
}

var RelationshipModel = Model{
	ID:     MakeModelID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name:   "relationship",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "leftName",
			ID:       MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		},
		Attribute{
			Name:     "rightName",
			ID:       MakeID("8d8524ab-92d6-49a3-8038-3ad957c5f6e8"),
			Datatype: String,
		},
		Attribute{
			Name:     "leftBinding",
			ID:       MakeID("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Int,
		},
		Attribute{
			Name:     "rightBinding",
			ID:       MakeID("4135be16-7c61-4750-b53d-f1eeff69086e"),
			Datatype: Int,
		},
	},
	RightRelationships: []Relationship{
		ModelRelationshipsRight,
		ModelRelationshipsLeft,
	},
}

var DatatypeModel = Model{
	ID:     MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name:   "datatype",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		},
		Attribute{
			Name:     "storedAs",
			ID:       MakeID("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: StoredAs,
		},
		Attribute{
			Name:     "enum",
			ID:       MakeID("931050f5-0022-4be2-87fb-d69537877a87"),
			Datatype: Bool,
		},
		Attribute{
			Name:     "native",
			ID:       MakeID("db56571e-1939-45f1-b122-9ecb8ad9fd7e"),
			Datatype: Bool,
		},
	},
	RightRelationships: []Relationship{
		AttributeDatatype,
	},
	LeftRelationships: []Relationship{
		ValidatorCode,
		DatatypeEnumValues,
	},
}

var CodeModel = Model{
	ID:     MakeModelID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name:   "code",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
			Datatype: String,
		},
		Attribute{
			Name:     "runtime",
			ID:       MakeID("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Runtime,
		},
		Attribute{
			Name:     "functionSignature",
			ID:       MakeID("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
			Datatype: FunctionSignature,
		},
		Attribute{
			Name:     "code",
			ID:       MakeID("80b3055b-08ad-41fe-b562-4a493bb6db36"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		ValidatorCode,
	},
}

var EnumValueModel = Model{
	ID:     MakeModelID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"),
	Name:   "enumValue",
	System: true,
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("5803e350-48f8-448d-9901-7c80f45c775b"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		DatatypeEnumValues,
	},
}

var ModelAttributes = Relationship{
	ID:           MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	LeftModelID:  MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "attributes",
	LeftBinding:  HasMany,
	RightModelID: MakeModelID("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	RightName:    "model",
	RightBinding: BelongsTo,
}

var ModelRelationshipsLeft = Relationship{
	ID:           MakeID("806334bf-98ce-4c08-87f4-5d9bed4f6d60"),
	LeftModelID:  MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "leftRelationships",
	LeftBinding:  HasMany,
	RightModelID: MakeModelID("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "leftModel",
	RightBinding: BelongsTo,
}

var ModelRelationshipsRight = Relationship{
	ID:           MakeID("3ccee9ea-f5b7-4707-9b0f-9e6f96ddd42e"),
	LeftModelID:  MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "rightRelationships",
	LeftBinding:  HasMany,
	RightModelID: MakeModelID("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "rightModel",
	RightBinding: BelongsTo,
}

var AttributeDatatype = Relationship{
	ID:           MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	LeftModelID:  MakeModelID("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	LeftName:     "datatype",
	LeftBinding:  BelongsTo,
	RightModelID: MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	RightName:    "attribute",
	RightBinding: HasOne,
}

var ValidatorCode = Relationship{
	ID:           MakeID("353a1d40-d292-47f6-b45c-06b059bed882"),
	LeftModelID:  MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "validator",
	LeftBinding:  BelongsTo,
	RightModelID: MakeModelID("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var DatatypeEnumValues = Relationship{
	ID:           MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
	LeftModelID:  MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "enumValues",
	LeftBinding:  HasMany,
	RightModelID: MakeModelID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"), // enum
	RightName:    "datatype",
	RightBinding: BelongsTo,
}

var boolValidator = Code{
	Name:              "bool",
	ID:                MakeID("8e806967-c462-47af-8756-48674537a909"),
	Runtime:           Native,
	Function:          datatypes.BoolFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var intValidator = Code{
	Name:              "int",
	ID:                MakeID("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Runtime:           Native,
	Function:          datatypes.IntFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var stringValidator = Code{
	Name:              "string",
	ID:                MakeID("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Runtime:           Native,
	Function:          datatypes.StringFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var textValidator = Code{
	Name:              "longText",
	ID:                MakeID("eb78bcc8-d933-46cb-808f-216a7042f7d0"),
	Runtime:           Native,
	Function:          datatypes.StringFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var uuidValidator = Code{
	Name:              "uuid",
	ID:                MakeID("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Runtime:           Native,
	Function:          datatypes.UUIDFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var floatValidator = Code{
	Name:              "float",
	ID:                MakeID("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Runtime:           Native,
	Function:          datatypes.FloatFromJSON,
	Executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var Bool = coreDatatype{
	ID:        MakeID("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:      "bool",
	Validator: boolValidator,
	StoredAs:  BoolStorage,
}

var Int = coreDatatype{
	ID:        MakeID("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:      "int",
	Validator: intValidator,
	StoredAs:  IntStorage,
}

var String = coreDatatype{
	ID:        MakeID("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:      "string",
	Validator: stringValidator,
	StoredAs:  StringStorage,
}

var LongText = coreDatatype{
	ID:        MakeID("1557f632-63b8-4a4a-86db-4ff7dc6f4457"),
	Name:      "longText",
	Validator: textValidator,
	StoredAs:  StringStorage,
}

var UUID = coreDatatype{
	ID:        MakeID("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:      "uuid",
	Validator: uuidValidator,
	StoredAs:  UUIDStorage,
}

var Float = coreDatatype{
	ID:        MakeID("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:      "float",
	Validator: floatValidator,
	StoredAs:  FloatStorage,
}

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
		ID:       MakeID("cecf8eac-d3be-4ca0-927a-127763d465b1"),
		Name:     "native",
		Datatype: MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"), //Runtime ID
	},
}

var Starlark = RuntimeEnumValue{
	EnumValue{
		ID:       MakeID("c0036590-8227-46cb-8cf9-689dd17616a3"),
		Name:     "starlark",
		Datatype: MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"), //Runtime ID
	},
}

var FromJSON = FunctionSignatureEnumValue{
	EnumValue{
		ID:       MakeID("508ba2cc-ce86-4615-bc0d-fe0d085a2051"),
		Name:     "fromJson",
		Datatype: MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"), //FunctionSignature ID
	},
}

var RPC = FunctionSignatureEnumValue{
	EnumValue{
		ID:       MakeID("8decedba-555b-47ca-a232-68100fbbf756"),
		Name:     "rpc",
		Datatype: MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"), //FunctionSignature ID
	},
}

var BoolStorage = StorageEnumValue{
	EnumValue{
		ID:       MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
		Name:     "bool",
		Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var IntStorage = StorageEnumValue{
	EnumValue{
		ID:       MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
		Name:     "int",
		Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var StringStorage = StorageEnumValue{
	EnumValue{
		ID:       MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
		Name:     "string",
		Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var FloatStorage = StorageEnumValue{
	EnumValue{
		ID:       MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
		Name:     "float",
		Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}

var UUIDStorage = StorageEnumValue{
	EnumValue{
		ID:       MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
		Name:     "uuid",
		Datatype: MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"), //StoredAs ID
	},
}
