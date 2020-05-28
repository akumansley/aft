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
		"storageType": Attribute{
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
