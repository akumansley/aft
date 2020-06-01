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
		"storageType": Attribute{
			Id:       uuid.MustParse("523edf8d-6ea5-4745-8182-98165a75d4da"),
			Datatype: NativeCode,
		},
		"jsonType": Attribute{
			Id:       uuid.MustParse("ad4cc765-de91-4b34-8e36-73031a190808"),
			Datatype: NativeCode,
		},
	},
	RightRelationships: []Relationship{
		AttributeDatatypes,
	},
	LeftRelationships: []Relationship{
		FromJsonCode,
		ToJsonCode,
	},
}

var CodeModel = Model{
	Id:   uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: map[string]Attribute{
		"name": Attribute{
			Id:       uuid.MustParse("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: String,
		},
		"runtime": Attribute{
			Id:       uuid.MustParse("e38e557c-7b18-4b8c-8be4-04ca7810c2c4"),
			Datatype: Enum,
		},
		"syntax": Attribute{
			Id:       uuid.MustParse("9ad0482e-92ab-45cd-b66b-24ddb1cc9971"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		FromJsonCode,
		ToJsonCode,
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

var AttributeDatatypes = Relationship{
	Id:           uuid.MustParse("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	LeftModelId:  uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"), // attribute
	LeftName:     "datatype",
	LeftBinding:  BelongsTo,
	RightModelId: uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	RightName:    "attribute",
	RightBinding: HasMany,
}

var FromJsonCode = Relationship{
	Id:           uuid.MustParse("353a1d40-d292-47f6-b45c-06b059bed882"),
	LeftModelId:  uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "fromJson",
	LeftBinding:  BelongsTo,
	RightModelId: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var ToJsonCode = Relationship{
	Id:           uuid.MustParse("3a7ee5c2-f93b-44bd-9f9d-58bd6ee592d7"),
	LeftModelId:  uuid.MustParse("c2ea9d6f-26ca-4674-b2b4-3a2bc3861a6a"), // datatype
	LeftName:     "toJson",
	LeftBinding:  BelongsTo,
	RightModelId: uuid.MustParse("8deaec0c-f281-4583-baf7-89c3b3b051f3"), // code
	RightName:    "datatype",
	RightBinding: HasOne,
}

var Bool = Datatype{
	Id:          uuid.MustParse("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:        "bool",
	FromJson:    boolFromJson,
	ToJson:      boolToJson,
	StorageType: false,
	JsonType:    false,
}

var Int = Datatype{
	Id:          uuid.MustParse("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:        "int",
	FromJson:    intFromJson,
	ToJson:      intToJson,
	StorageType: int64(0),
	JsonType:    0.0,
}

var Enum = Datatype{
	Id:          uuid.MustParse("f9e66ef9-2fa3-4588-81c1-b7be6a28352e"),
	Name:        "enum",
	FromJson:    enumFromJson,
	ToJson:      enumToJson,
	StorageType: int64(0),
	JsonType:    0.0,
}

var String = Datatype{
	Id:          uuid.MustParse("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:        "string",
	FromJson:    stringFromJson,
	ToJson:      stringToJson,
	StorageType: "",
	JsonType:    "",
}

var Text = Datatype{
	Id:          uuid.MustParse("4b601851-421d-4633-8a68-7fefea041361"),
	Name:        "text",
	FromJson:    textFromJson,
	ToJson:      textToJson,
	StorageType: "",
	JsonType:    "",
}

var EmailAddress = Datatype{
	Id:          uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:        "emailAddress",
	FromJson:    emailAddressFromJson,
	ToJson:      emailAddressToJson,
	StorageType: "",
	JsonType:    "",
}

var UUID = Datatype{
	Id:          uuid.MustParse("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:        "uuid",
	FromJson:    uuidFromJson,
	ToJson:      uuidToJson,
	StorageType: uuid.UUID{},
	JsonType:    uuid.UUID{},
}

var Float = Datatype{
	Id:          uuid.MustParse("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:        "float",
	FromJson:    floatFromJson,
	ToJson:      floatToJson,
	StorageType: 0.0,
	JsonType:    0.0,
}

var URL = Datatype{
	Id:          uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:        "url",
	FromJson:    urlFromJson,
	ToJson:      urlToJson,
	StorageType: "",
	JsonType:    "",
}

var NativeCode = Datatype{
	Id:          uuid.MustParse("f34e5dd5-9209-4ce0-81ef-8e2d1ee86ece"),
	Name:        "nativeCode",
	FromJson:    nativeCodeFromJson,
	ToJson:      nativeCodeToJson,
	StorageType: "",
	JsonType:    "",
}

var boolFromJson = Code{
	Id:      uuid.MustParse("8e806967-c462-47af-8756-48674537a909"),
	Name:    "boolFromJson",
	Runtime: Golang,
}

var boolToJson = Code{
	Id:      uuid.MustParse("22bb89d5-1656-4a31-9458-95c133a3abc3"),
	Name:    "boolToJson",
	Runtime: Golang,
}

var intFromJson = Code{
	Id:      uuid.MustParse("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	Name:    "intFromJson",
	Runtime: Golang,
}

var intToJson = Code{
	Id:      uuid.MustParse("21120409-dd95-479e-9aa2-01d01418e65f"),
	Name:    "intToJson",
	Runtime: Golang,
}

var enumFromJson = Code{
	Id:      uuid.MustParse("5c3b9da9-c592-41da-b6e2-8c8dd97186c3"),
	Name:    "enumFromJson",
	Runtime: Golang,
}

var enumToJson = Code{
	Id:      uuid.MustParse("367acb04-69d1-492b-953e-b26488f10390"),
	Name:    "enumToJson",
	Runtime: Golang,
}

var stringFromJson = Code{
	Id:      uuid.MustParse("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	Name:    "stringFromJson",
	Runtime: Golang,
}

var stringToJson = Code{
	Id:      uuid.MustParse("a0f3f396-5ce4-4b12-ad92-39bb1df2d1cb"),
	Name:    "stringToJson",
	Runtime: Golang,
}

var textFromJson = Code{
	Id:      uuid.MustParse("9f10ac9f-afd2-423a-8857-d900a0c97563"),
	Name:    "textFromJson",
	Runtime: Golang,
}

var textToJson = Code{
	Id:      uuid.MustParse("0fa33363-1dd0-4898-963b-fce064144cef"),
	Name:    "textToJson",
	Runtime: Golang,
}

var emailAddressFromJson = Code{
	Id:      uuid.MustParse("ed046b08-ade2-4570-ade4-dd1e31078219"),
	Name:    "emailAddressFromJson",
	Runtime: Golang,
}

var emailAddressToJson = Code{
	Id:      uuid.MustParse("6a26a584-5198-40fc-82cb-1225411fbafb"),
	Name:    "emailAddressToJson",
	Runtime: Golang,
}

var uuidFromJson = Code{
	Id:      uuid.MustParse("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	Name:    "uuidFromJson",
	Runtime: Golang,
}

var uuidToJson = Code{
	Id:      uuid.MustParse("810fbb58-25d6-4ccf-a451-0f5fc543fa5d"),
	Name:    "uuidToJson",
	Runtime: Golang,
}

var floatFromJson = Code{
	Id:      uuid.MustParse("83a5f999-00b0-4bc1-879a-434869cf7301"),
	Name:    "floatFromJson",
	Runtime: Golang,
}

var floatToJson = Code{
	Id:      uuid.MustParse("b18aa08a-2080-4d6c-bd3e-df93d62d80cc"),
	Name:    "floatToJson",
	Runtime: Golang,
}

var urlFromJson = Code{
	Id:      uuid.MustParse("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	Name:    "urlFromJson",
	Runtime: Golang,
}

var urlToJson = Code{
	Id:      uuid.MustParse("4f61e364-1fcb-4099-b13d-1dec1fb14f9a"),
	Name:    "urlToJson",
	Runtime: Golang,
}

var nativeCodeFromJson = Code{
	Id:      uuid.MustParse("05f76f51-8805-4c7d-8087-6e3315a64807"),
	Name:    "nativeCodeFromJson",
	Runtime: Golang,
}

var nativeCodeToJson = Code{
	Id:      uuid.MustParse("940cd1ba-0230-4210-9b0a-791d1b642be5"),
	Name:    "nativeCodeToJson",
	Runtime: Golang,
}
