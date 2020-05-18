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
	Id:   uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "user",
	Attributes: map[string]Attribute{
		"firstName": Attribute{
			Id:       uuid.MustParse("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			AttrType: String,
		},
		"lastName": Attribute{
			Id:       uuid.MustParse("462212e7-dd94-403e-8314-e271fd7ccec9"),
			AttrType: String,
		},
		"age": Attribute{
			Id:       uuid.MustParse("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			AttrType: Int,
		},
	},
	LeftRelationships: []Relationship{
		UserPosts,
		UserProfile,
	},
}

var UserPosts = Relationship{
	Id:           uuid.MustParse("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
	LeftModelId:  uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"), // user
	LeftName:     "posts",
	LeftBinding:  HasMany,
	RightModelId: uuid.MustParse("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"), // post
	RightName:    "author",
	RightBinding: BelongsTo,
}

var UserProfile = Relationship{
	Id:           uuid.MustParse("52a31e61-f1d3-4091-8dcf-78236ef84f6f"),
	LeftModelId:  uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"), // user
	LeftName:     "profile",
	LeftBinding:  HasOne,
	RightModelId: uuid.MustParse("66783192-4111-4bd8-95dd-e7da460378df"), // profile
	RightName:    "user",
	RightBinding: BelongsTo,
}

var Profile = Model{
	Id:   uuid.MustParse("66783192-4111-4bd8-95dd-e7da460378df"),
	Name: "profile",
	Attributes: map[string]Attribute{
		"text": Attribute{
			Id:       uuid.MustParse("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			AttrType: String,
		},
	},
	RightRelationships: []Relationship{
		UserProfile,
	},
}

var Post = Model{
	Id:   uuid.MustParse("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	Name: "post",
	Attributes: map[string]Attribute{
		"text": Attribute{
			Id:       uuid.MustParse("b3af6694-b621-43a2-be7f-00956fa505c0"),
			AttrType: String,
		},
	},
	RightRelationships: []Relationship{
		UserPosts,
	},
}
