package db

import (
	"github.com/google/uuid"
)

func AddSampleModels(db DB) {
	tx := db.NewRWTx()
	tx.SaveModel(User)
	tx.SaveModel(Profile)
	tx.SaveModel(Post)
	tx.Commit()
}

var User = Model{
	Type: "model",
	Id:   uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "user",
	Attributes: map[string]Attribute{
		"firstName": Attribute{
			Id:       uuid.MustParse("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			Type:     "attribute",
			AttrType: String,
		},
		"lastName": Attribute{
			Id:       uuid.MustParse("462212e7-dd94-403e-8314-e271fd7ccec9"),
			Type:     "attribute",
			AttrType: String,
		},
		"age": Attribute{
			Id:       uuid.MustParse("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			Type:     "attribute",
			AttrType: Int,
		},
	},
	Relationships: map[string]Relationship{
		"posts": Relationship{
			Id:          uuid.MustParse("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
			Type:        "relationship",
			TargetModel: "post",
			TargetRel:   "author",
			RelType:     HasMany,
		},
		"profile": Relationship{
			Id:          uuid.MustParse("c4043a82-a3df-4d55-ac76-c8412131d34a"),
			Type:        "relationship",
			TargetModel: "profile",
			TargetRel:   "user",
			RelType:     HasOne,
		},
	},
}

var Profile = Model{
	Type: "model",
	Id:   uuid.MustParse("66783192-4111-4bd8-95dd-e7da460378df"),
	Name: "profile",
	Attributes: map[string]Attribute{
		"text": Attribute{
			Id:       uuid.MustParse("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			Type:     "attribute",
			AttrType: String,
		},
	},
	Relationships: map[string]Relationship{
		"user": Relationship{
			Id:          uuid.MustParse("c3172b78-e091-4040-b686-a0a5a844117a"),
			Type:        "relationship",
			TargetModel: "user",
			TargetRel:   "profile",
			RelType:     BelongsTo,
		},
	},
}

var Post = Model{
	Type: "model",
	Id:   uuid.MustParse("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	Name: "post",
	Attributes: map[string]Attribute{
		"text": Attribute{
			Id:       uuid.MustParse("b3af6694-b621-43a2-be7f-00956fa505c0"),
			Type:     "attribute",
			AttrType: String,
		},
	},
	Relationships: map[string]Relationship{
		"author": Relationship{
			Id:          uuid.MustParse("0ea3b703-8d6c-4aa6-8aa0-68bc3fd39eb0"),
			Type:        "relationship",
			TargetModel: "user",
			TargetRel:   "posts",
			RelType:     BelongsTo,
		},
	},
}
