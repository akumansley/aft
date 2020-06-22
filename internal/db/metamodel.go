package db

import (
	"awans.org/aft/internal/datatypes"
)

var ModelModel = Model{
	ID:   MakeModelID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Datatype: String,
		},
	},
	LeftRelationships: []Relationship{
		ModelAttributes,
		ModelRelationshipsLeft,
		ModelRelationshipsRight,
	},
}

var AttributeModel = Model{
	ID:   MakeModelID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
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
	ID:   MakeModelID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
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
	ID:   MakeModelID("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name: "datatype",
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		},
		Attribute{
			Name:     "storedAs",
			ID:       MakeID("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: Enum,
		},
	},
	RightRelationships: []Relationship{
		AttributeDatatype,
	},
	LeftRelationships: []Relationship{
		ValidatorCode,
	},
}

var CodeModel = Model{
	ID:   MakeModelID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: []Attribute{
		Attribute{
			Name:     "name",
			ID:       MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
			Datatype: String,
		},
		Attribute{
			Name:     "runtime",
			ID:       MakeID("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Enum,
		},
		Attribute{
			Name:     "functionSignature",
			ID:       MakeID("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
			Datatype: Enum,
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

var boolValidator = Code{
	Name:              "bool",
	ID:                MakeID("8e806967-c462-47af-8756-48674537a909"),
	Runtime:           Golang,
	Function:          datatypes.BoolFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var intValidator = Code{
	Name:              "int",
	ID:                MakeID("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Runtime:           Golang,
	Function:          datatypes.IntFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var enumValidator = Code{
	Name:              "enum",
	ID:                MakeID("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
	Runtime:           Golang,
	Function:          datatypes.EnumFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var stringValidator = Code{
	Name:              "string",
	ID:                MakeID("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Runtime:           Golang,
	Function:          datatypes.StringFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var textValidator = Code{
	Name:              "text",
	ID:                MakeID("9f10ac9f-afd2-423a-8857-d900a0c97563"),
	Runtime:           Golang,
	Function:          datatypes.TextFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var uuidValidator = Code{
	Name:              "uuid",
	ID:                MakeID("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Runtime:           Golang,
	Function:          datatypes.UUIDFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var floatValidator = Code{
	Name:              "float",
	ID:                MakeID("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Runtime:           Golang,
	Function:          datatypes.FloatFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var Bool = Datatype{
	ID:        MakeID("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:      "bool",
	Validator: boolValidator,
	StoredAs:  BoolStorage,
}

var Int = Datatype{
	ID:        MakeID("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:      "int",
	Validator: intValidator,
	StoredAs:  IntStorage,
}

var Enum = Datatype{
	ID:        MakeID("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:      "enum",
	Validator: enumValidator,
	StoredAs:  IntStorage,
}

var String = Datatype{
	ID:        MakeID("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:      "string",
	Validator: stringValidator,
	StoredAs:  StringStorage,
}

var Text = Datatype{
	ID:        MakeID("4b601851-421d-4633-8a68-7fefea041361"),
	Name:      "text",
	Validator: textValidator,
	StoredAs:  StringStorage,
}

var UUID = Datatype{
	ID:        MakeID("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:      "uuid",
	Validator: uuidValidator,
	StoredAs:  UUIDStorage,
}

var Float = Datatype{
	ID:        MakeID("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:      "float",
	Validator: floatValidator,
	StoredAs:  FloatStorage,
}
