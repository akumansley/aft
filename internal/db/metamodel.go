package db

var ModelModel = ModelL{
	ID:   MakeID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var ModelAttributeInterface = InterfaceL{
	ID:   MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "modelAttribute",
	Attributes: []Attribute{
		InterfaceAttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsInterface()

// a concrete model for storing the attributes of interfaces
var InterfaceAttributeModel = ModelL{
	ID:   MakeID("41daafd2-b2cd-45c8-a087-84464b674a58"),
	Name: "interfaceAttribute",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("0391437a-cd60-449c-9689-57666535b9e6"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var ConcreteAttributeModel = ModelL{
	ID:   MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "concreteAttribute",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var ComputedAttrModel = ModelL{
	ID:   MakeID("29d88992-1855-4abc-968c-cf06e0979420"),
	Name: "computedAttribute",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("dc8ee712-f8a1-4b72-bbf7-f17d74beb796"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var RelationshipModel = ModelL{
	ID:   MakeID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		}.AsAttribute(),
		ConcreteAttributeL{
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
	Target: CoreDatatypeModel,
	Multi:  false,
}.AsRelationship()

var RelationshipSource = RelationshipL{
	Name:   "source",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: RelationshipModel,
	Target: ModelModel,
	Multi:  false,
}.AsRelationship()

var RelationshipTarget = RelationshipL{
	Name:   "target",
	ID:     MakeID("e194f9bf-ea7a-4c78-a179-bdf9c044ac3c"),
	Source: RelationshipModel,
	Target: ModelModel,
	Multi:  false,
}.AsRelationship()

var CoreDatatypeModel = ModelL{
	ID:   MakeID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name: "coreDatatype",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		}.AsAttribute(),
		ConcreteAttributeL{
			Name:     "storedAs",
			ID:       MakeID("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: StoredAs,
		}.AsAttribute(),
	},
}.AsModel()

var EnumValueModel = ModelL{
	ID:   MakeID("b0f2f6d1-9e7e-4ffe-992f-347b2d0731ac"),
	Name: "enumValue",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("5803e350-48f8-448d-9901-7c80f45c775b"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var DatatypeValidator = RelationshipL{
	ID:     MakeID("353a1d40-d292-47f6-b45c-06b059bed882"),
	Name:   "validator",
	Source: CoreDatatypeModel,   // datatype
	Target: NativeFunctionModel, // code
	Multi:  false,
}.AsRelationship()

// var DatatypeEnumValues = RelationshipL{
// 	ID:     MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
// 	Source: EnumDatatypeModel, // datatype
// 	Name:   "enumValues",
// 	Multi:  true,
// 	Target: EnumValueModel,
// }.AsRelationship()

// var enumValidator = NativeFunctionL{
// 	Name:              "enum",
// 	ID:                MakeID("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
// 	Runtime:           Native,
// 	Function:          datatypes.EnumFromJSON,
// 	FunctionSignature: FromJSON,
// }.AsFunction()

// var textValidator = NativeFunctionL{
// 	Name:              "text",
// 	ID:                MakeID("9f10ac9f-afd2-423a-8857-d900a0c97563"),
// 	Runtime:           Native,
// 	Function:          datatypes.TextFromJSON,
// 	FunctionSignature: FromJSON,
// }.AsFunction()

var FunctionSignature = EnumL{
	ID:   MakeID("45c261f8-b54a-4e78-9c3c-5383cb99fe20"),
	Name: "functionSignature",
	Values: []EnumValue{
		FromJSON,
		RPC,
	},
}.AsDatatype()

var StoredAs = EnumL{
	ID:   MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	Name: "storedAs",
	Values: []EnumValue{
		BoolStorage,
		IntStorage,
		StringStorage,
		FloatStorage,
		UUIDStorage,
	},
}.AsDatatype()

var FromJSON = EnumValueL{
	ID:   MakeID("508ba2cc-ce86-4615-bc0d-fe0d085a2051"),
	Name: "fromJson",
}.AsEnumValue()

var RPC = EnumValueL{
	ID:   MakeID("8decedba-555b-47ca-a232-68100fbbf756"),
	Name: "rpc",
}.AsEnumValue()

var BoolStorage = EnumValueL{
	ID:   MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
	Name: "bool",
}.AsEnumValue()

var IntStorage = EnumValueL{
	ID:   MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
	Name: "int",
}.AsEnumValue()

var StringStorage = EnumValueL{
	ID:   MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
	Name: "string",
}.AsEnumValue()

var FloatStorage = EnumValueL{
	ID:   MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
	Name: "float",
}.AsEnumValue()

var UUIDStorage = EnumValueL{
	ID:   MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
	Name: "uuid",
}.AsEnumValue()
