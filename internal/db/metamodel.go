package db

import (
	"github.com/google/uuid"
)

var ModelModel = Model{
	Id:   uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: map[string]Attribute{
		"name": Attribute{
			Id:       uuid.MustParse("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Datatype: String,
		},
	},
	LeftRelationships: []Relationship{
		ModelAttributes,
		ModelRelationshipsLeft,
		ModelRelationshipsRight,
	},
}

var ModelAttributes = Relationship{
	Id:           uuid.MustParse("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	LeftModelId:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "attributes",
	LeftBinding:  HasMany,
	RightModelId: uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	RightName:    "model",
	RightBinding: BelongsTo,
}

var ModelRelationshipsLeft = Relationship{
	Id:           uuid.MustParse("806334bf-98ce-4c08-87f4-5d9bed4f6d60"),
	LeftModelId:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "leftRelationships",
	LeftBinding:  HasMany,
	RightModelId: uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "leftModel",
	RightBinding: BelongsTo,
}

var ModelRelationshipsRight = Relationship{
	Id:           uuid.MustParse("3ccee9ea-f5b7-4707-9b0f-9e6f96ddd42e"),
	LeftModelId:  uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"), // model
	LeftName:     "rightRelationships",
	LeftBinding:  HasMany,
	RightModelId: uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"), // relationship
	RightName:    "rightModel",
	RightBinding: BelongsTo,
}

var AttributeModel = Model{
	Id:   uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
	Attributes: map[string]Attribute{
		"name": Attribute{
			Id:       uuid.MustParse("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
		"datatype": Attribute{
			Id:       uuid.MustParse("c29a6558-7676-40a8-be00-e0933342efd7"),
			Datatype: UUID,
		},
	},
	RightRelationships: []Relationship{
		ModelAttributes,
	},
	LeftRelationships: []Relationship{
		AttributeDatatypes,
	},
}

var RelationshipModel = Model{
	Id:   uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: map[string]Attribute{
		"leftName": Attribute{
			Id:       uuid.MustParse("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Datatype: String,
		},
		"rightName": Attribute{
			Id:       uuid.MustParse("8d8524ab-92d6-49a3-8038-3ad957c5f6e8"),
			Datatype: String,
		},
		"leftBinding": Attribute{
			Id:       uuid.MustParse("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Datatype: Int,
		},
		"rightBinding": Attribute{
			Id:       uuid.MustParse("4135be16-7c61-4750-b53d-f1eeff69086e"),
			Datatype: Int,
		},
	},
	RightRelationships: []Relationship{
		ModelRelationshipsRight,
		ModelRelationshipsLeft,
	},
}

var DatatypeModel = Model{
	Id:   uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"),
	Name: "datatype",
	Attributes: map[string]Attribute{
		"name": Attribute{
			Id:       uuid.MustParse("0a0fe2bc-7443-4111-8b49-9fe41f186261"),
			Datatype: String,
		},
		"fromJson": Attribute{
			Id:       uuid.MustParse("ebe07b17-8c2c-4214-b97f-9f833059ed4e"),
			Datatype: NativeCode,
		},
		"type": Attribute{
			Id:       uuid.MustParse("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		AttributeDatatypes,
	},
}

var AttributeDatatypes = Relationship{
	Id:           uuid.MustParse("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	LeftModelId:  uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	LeftName:     "datatype",
	LeftBinding:  BelongsTo,
	RightModelId: uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), //Datatype
	RightName:    "attribute",
	RightBinding: HasMany,
}

var Bool = Datatype{
	Id:       uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:     "bool",
	FromJson: boolFromJson,
	Type:     false,
}

var Int = Datatype{
	Id:       uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:     "int",
	FromJson: intFromJson,
	Type:     int64(0),
}

var Enum = Datatype{
	Id:       uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:     "enum",
	FromJson: enumFromJson,
	Type:     int64(0),
}

var String = Datatype{
	Id:       uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:     "string",
	FromJson: stringFromJson,
	Type:     "",
}

var Text = Datatype{
	Id:       uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:     "text",
	FromJson: textFromJson,
	Type:     "",
}

var EmailAddress = Datatype{
	Id:       uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:     "email address",
	FromJson: emailAddressFromJson,
	Type:     "",
}

var UUID = Datatype{
	Id:       uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:     "uuid",
	FromJson: uuidFromJson,
	Type:     uuid.UUID{},
}
var Float = Datatype{
	Id:       uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:     "float",
	FromJson: floatFromJson,
	Type:     0.0,
}

var URL = Datatype{
	Id:       uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:     "url",
	FromJson: urlFromJson,
	Type:     "",
}

var NativeCode = Datatype{
	Id:       uuid.MustParse("f34e5dd5-9209-4ce0-81ef-8e2d1ee86ece"),
	Name:     "native code",
	FromJson: nativeCodeFromJson,
	Type:     "",
}

var Javascript = Datatype{
	Id:       uuid.MustParse("210fd74f-bb52-4ff6-b45a-d9fcd8d980ae"),
	Name:     "javascript",
	FromJson: javascriptFromJson,
	Type:     "",
}
