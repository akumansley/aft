package db

import (
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
		"type": Attribute{
			ID:       uuid.MustParse("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: Enum,
		},
	},
	RightRelationships: []Relationship{
		AttributeDatatype,
	},
	LeftRelationships: []Relationship{
		FromJSONCode,
		ToJSONCode,
	},
}

var CodeModel = Model{
	ID:   uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: map[string]Attribute{
		"name": Attribute{
			ID:       uuid.MustParse("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: String,
		},
		"function": Attribute{
			ID:       uuid.MustParse("9ad0482e-92ab-45cd-b66b-24ddb1cc9971"),
			Datatype: Enum,
		},
		"runtime": Attribute{
			ID:       uuid.MustParse("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Enum,
		},
		"code": Attribute{
			ID:       uuid.MustParse("9ad0482e-92ab-45cd-b66b-24ddb1cc9971"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		FromJSONCode,
		ToJSONCode,
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

var FromJSONCode = Relationship{
	ID:           uuid.MustParse("353a1d40-d292-47f6-b45c-06b059bed882"),
	LeftModelID:  uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "fromJson",
	LeftBinding:  BelongsTo,
	RightModelID: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var ToJSONCode = Relationship{
	ID:           uuid.MustParse("3a7ee5c2-f93b-44bd-9f9d-58bd6ee592d7"),
	LeftModelID:  uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "toJson",
	LeftBinding:  BelongsTo,
	RightModelID: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var Bool = Datatype{
	ID:       uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:     "bool",
	FromJSON: boolFromJSON,
	ToJSON:   boolToJSON,
	Type:     BoolType,
}

var Int = Datatype{
	ID:       uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:     "int",
	FromJSON: intFromJSON,
	ToJSON:   intToJSON,
	Type:     IntType,
}

var Enum = Datatype{
	ID:       uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:     "enum",
	FromJSON: enumFromJSON,
	ToJSON:   enumToJSON,
	Type:     IntType,
}

var String = Datatype{
	ID:       uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:     "string",
	FromJSON: stringFromJSON,
	ToJSON:   stringToJSON,
	Type:     StringType,
}

var Text = Datatype{
	ID:       uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:     "text",
	FromJSON: textFromJSON,
	ToJSON:   textToJSON,
	Type:     StringType,
}

var EmailAddress = Datatype{
	ID:       uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:     "emailAddress",
	FromJSON: emailAddressFromJSON,
	ToJSON:   emailAddressToJSON,
	Type:     StringType,
}

var UUID = Datatype{
	ID:       uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:     "uuid",
	FromJSON: uuidFromJSON,
	ToJSON:   uuidToJSON,
	Type:     UUIDType,
}

var Float = Datatype{
	ID:       uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:     "float",
	FromJSON: floatFromJSON,
	ToJSON:   floatToJSON,
	Type:     FloatType,
}

var URL = Datatype{
	ID:       uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:     "url",
	FromJSON: URLFromJSON,
	ToJSON:   URLToJSON,
	Type:     StringType,
}

var boolFromJSON = Code{
	ID:       uuid.MustParse("8e806967-c462-47af-8756-48674537a909"),
	Name:     "boolFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var boolToJSON = Code{
	ID:       uuid.MustParse("22bb89d5-1656-4a31-9458-95c133a3abc3"),
	Name:     "boolToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var intFromJSON = Code{
	ID:       uuid.MustParse("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Name:     "intFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var intToJSON = Code{
	ID:       uuid.MustParse("21120409-dd95-479e-9aa2-01d01418e65f"),
	Name:     "intToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var enumFromJSON = Code{
	ID:       uuid.MustParse("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
	Name:     "enumFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var enumToJSON = Code{
	ID:       uuid.MustParse("367acb04-69d1-492b-953e-b26488f10390"),
	Name:     "enumToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var stringFromJSON = Code{
	ID:       uuid.MustParse("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Name:     "stringFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var stringToJSON = Code{
	ID:       uuid.MustParse("a0f3f396-5ce4-4b12-ad92-39bb1df2d1cb"),
	Name:     "stringToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var textFromJSON = Code{
	ID:       uuid.MustParse("9f10ac9f-afd2-423a-8857-d900a0c97563"),
	Name:     "textFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var textToJSON = Code{
	ID:       uuid.MustParse("0fa33363-1dd0-4898-963b-fce064144cef"),
	Name:     "textToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var emailAddressFromJSON = Code{
	ID:       uuid.MustParse("ed046b08-ade2-4570-ade4-dd1e31078219"),
	Name:     "emailAddressFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var emailAddressToJSON = Code{
	ID:       uuid.MustParse("6a26a584-5198-40fc-82cb-1225411fbafb"),
	Name:     "emailAddressToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var uuidFromJSON = Code{
	ID:       uuid.MustParse("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Name:     "uuidFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var uuidToJSON = Code{
	ID:       uuid.MustParse("810fbb58-25d6-4ccf-a451-0f5fc543fa5d"),
	Name:     "uuidToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var floatFromJSON = Code{
	ID:       uuid.MustParse("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Name:     "floatFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var floatToJSON = Code{
	ID:       uuid.MustParse("b18aa08a-2080-4d6c-bd3e-df93d62d80cc"),
	Name:     "floatToJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var URLFromJSON = Code{
	ID:       uuid.MustParse("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	Name:     "urlFromJson",
	Runtime:  Golang,
	Function: FromJSON,
}

var URLToJSON = Code{
	ID:       uuid.MustParse("4f61e364-1fcb-4099-b13d-1dec1fb14f9a"),
	Name:     "urlToJson",
	Runtime:  Golang,
	Function: FromJSON,
}
