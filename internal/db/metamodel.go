package db

import (
	"awans.org/aft/internal/datatypes"
	"github.com/google/uuid"
)

var ModelModel = Model{
	ID:   uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: map[string]Attribute{
		"name": Attribute{
			ID:       uuid.MustParse("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
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
	ID:   uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
	Attributes: map[string]Attribute{
		"name": Attribute{
			ID:       uuid.MustParse("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
		"datatypeId": Attribute{ //todo remove hack
			ID:       uuid.MustParse("bfeefcbf-b9f7-44e6-9951-134755f7e1cd"),
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
	ID:   uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: map[string]Attribute{
		"leftName": Attribute{
			ID:       uuid.MustParse("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		},
		"rightName": Attribute{
			ID:       uuid.MustParse("8d8524ab-92d6-49a3-8038-3ad957c5f6e8"),
			Datatype: String,
		},
		"leftBinding": Attribute{
			ID:       uuid.MustParse("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Int,
		},
		"rightBinding": Attribute{
			ID:       uuid.MustParse("4135be16-7c61-4750-b53d-f1eeff69086e"),
			Datatype: Int,
		},
	},
	RightRelationships: []Relationship{
		ModelRelationshipsRight,
		ModelRelationshipsLeft,
	},
}

var DatatypeModel = Model{
	ID:   uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name: "datatype",
	Attributes: map[string]Attribute{
		"name": Attribute{
			ID:       uuid.MustParse("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		},
		"storedAs": Attribute{
			ID:       uuid.MustParse("523edf8d-6ea5-4745-8182-98165a75d4da"),
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
	ID:   uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: map[string]Attribute{
		"name": Attribute{
			ID:       uuid.MustParse("c47bcd30-01ea-467f-ad02-114342070241"),
			Datatype: String,
		},
		"runtime": Attribute{
			ID:       uuid.MustParse("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Enum,
		},
		"functionSignature": Attribute{
			ID:       uuid.MustParse("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
			Datatype: Enum,
		},
		"code": Attribute{
			ID:       uuid.MustParse("80b3055b-08ad-41fe-b562-4a493bb6db36"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		ValidatorCode,
	},
}

var ModelAttributes = Relationship{
	ID:           uuid.MustParse("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	LeftModelID:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "attributes",
	LeftBinding:  HasMany,
	RightModelID: uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	RightName:    "model",
	RightBinding: BelongsTo,
}

var ModelRelationshipsLeft = Relationship{
	ID:           uuid.MustParse("806334bf-98ce-4c08-87f4-5d9bed4f6d60"),
	LeftModelID:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "leftRelationships",
	LeftBinding:  HasMany,
	RightModelID: uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "leftModel",
	RightBinding: BelongsTo,
}

var ModelRelationshipsRight = Relationship{
	ID:           uuid.MustParse("3ccee9ea-f5b7-4707-9b0f-9e6f96ddd42e"),
	LeftModelID:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "rightRelationships",
	LeftBinding:  HasMany,
	RightModelID: uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "rightModel",
	RightBinding: BelongsTo,
}

var AttributeDatatype = Relationship{
	ID:           uuid.MustParse("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	LeftModelID:  uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	LeftName:     "datatype",
	LeftBinding:  BelongsTo,
	RightModelID: uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	RightName:    "attribute",
	RightBinding: HasOne,
}

var ValidatorCode = Relationship{
	ID:           uuid.MustParse("353a1d40-d292-47f6-b45c-06b059bed882"),
	LeftModelID:  uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "validator",
	LeftBinding:  BelongsTo,
	RightModelID: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var boolValidator = Code{
	ID:                uuid.MustParse("8e806967-c462-47af-8756-48674537a909"),
	Name:              "boolValidator",
	Runtime:           Golang,
	Function:          datatypes.BoolFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var intValidator = Code{
	ID:                uuid.MustParse("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Name:              "intValidator",
	Runtime:           Golang,
	Function:          datatypes.IntFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var enumValidator = Code{
	ID:                uuid.MustParse("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
	Name:              "enumValidator",
	Runtime:           Golang,
	Function:          datatypes.EnumFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var stringValidator = Code{
	ID:                uuid.MustParse("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Name:              "stringValidator",
	Runtime:           Golang,
	Function:          datatypes.StringFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var textValidator = Code{
	ID:                uuid.MustParse("9f10ac9f-afd2-423a-8857-d900a0c97563"),
	Name:              "textValidator",
	Runtime:           Golang,
	Function:          datatypes.TextFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var uuidValidator = Code{
	ID:                uuid.MustParse("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Name:              "uuidValidator",
	Runtime:           Golang,
	Function:          datatypes.UUIDFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var floatValidator = Code{
	ID:                uuid.MustParse("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Name:              "floatValidator",
	Runtime:           Golang,
	Function:          datatypes.FloatFromJSON,
	executor:          &bootstrapCodeExecutor{},
	FunctionSignature: FromJSON,
}

var Bool = Datatype{
	ID:        uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:      "bool",
	Validator: boolValidator,
	StoredAs:  BoolStorage,
}

var Int = Datatype{
	ID:        uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:      "int",
	Validator: intValidator,
	StoredAs:  IntStorage,
}

var Enum = Datatype{
	ID:        uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:      "enum",
	Validator: enumValidator,
	StoredAs:  IntStorage,
}

var String = Datatype{
	ID:        uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:      "string",
	Validator: stringValidator,
	StoredAs:  StringStorage,
}

var Text = Datatype{
	ID:        uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:      "text",
	Validator: textValidator,
	StoredAs:  StringStorage,
}

var UUID = Datatype{
	ID:        uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:      "uuid",
	Validator: uuidValidator,
	StoredAs:  UUIDStorage,
}

var Float = Datatype{
	ID:        uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:      "float",
	Validator: floatValidator,
	StoredAs:  FloatStorage,
}
