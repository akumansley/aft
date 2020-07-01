package db

var ModelAttributeInterface = InterfaceL{
	ID:   MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "modelAttribute",
	Attributes: []InterfaceAttributeL{
		InterfaceAttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
	},
}

var InterfaceModel = ModelL{
	ID:   MakeID("7a16a48d-8827-4e70-b982-d85af04c4ec9"),
	Name: "interface",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("cb9001df-b8d2-467c-87da-196057c74946"),
			Datatype: String,
		},
	},
}

// a concrete model for storing the attributes of interfaces
var InterfaceAttributeModel = ModelL{
	ID:   MakeID("41daafd2-b2cd-45c8-a087-84464b674a58"),
	Name: "interfaceAttribute",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("0391437a-cd60-449c-9689-57666535b9e6"),
			Datatype: String,
		},
	},
}

var RelationshipModel = ModelL{
	ID:   MakeID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		},
		ConcreteAttributeL{
			Name:     "multi",
			ID:       MakeID("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Bool,
		},
	},
}

var ModelAttributes = RelationshipL{
	Name:   "attributes",
	ID:     MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	Source: ModelModel,
	Target: ModelAttributeInterface,
	Multi:  true,
}

var AttributeDatatype = RelationshipL{
	Name:   "datatype",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Target: CoreDatatypeModel,
	Multi:  false,
}

var RelationshipSource = RelationshipL{
	Name:   "source",
	ID:     MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	Source: RelationshipModel,
	Target: ModelModel,
	Multi:  false,
}

var RelationshipTarget = RelationshipL{
	Name:   "target",
	ID:     MakeID("e194f9bf-ea7a-4c78-a179-bdf9c044ac3c"),
	Source: RelationshipModel,
	Target: ModelModel,
	Multi:  false,
}

var EnumModel = ModelL{
	ID:   MakeID(""),
	Name: "enum",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID(""),
			Datatype: String,
		},
	},
}

var EnumEnumValues = RelationshipL{
	ID:     MakeID("7f9aa1bc-dd19-4db9-9148-bf302c9d99da"),
	Name:   "enumValues",
	Source: EnumModel,
	Target: EnumValueModel,
	Multi:  true,
}

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
	Values: []EnumValueL{
		FromJSON,
		RPC,
		Getter,
		Setter,
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

var Getter = EnumValueL{
	ID:   MakeID("8ec4cd0e-72c5-4c75-9576-11202c0e562d"),
	Name: "getter",
}

var Setter = EnumValueL{
	ID:   MakeID("3623a700-0813-48d1-a14d-ef1bc2aa3503"),
	Name: "setter",
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
