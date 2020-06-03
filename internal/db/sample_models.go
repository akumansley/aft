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
	ID:   uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "user",
	Attributes: map[string]Attribute{
		"firstName": Attribute{
			ID:       uuid.MustParse("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			Datatype: Andrew,
		},
		"lastName": Attribute{
			ID:       uuid.MustParse("462212e7-dd94-403e-8314-e271fd7ccec9"),
			Datatype: String,
		},
		"age": Attribute{
			ID:       uuid.MustParse("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			Datatype: Int,
		},
		"emailAddress": Attribute{
			ID:       uuid.MustParse("0fe6bd01-9828-43ac-b004-37620083344d"),
			Datatype: EmailAddress,
		},
	},
	LeftRelationships: []Relationship{
		UserPosts,
		UserProfile,
	},
}

var UserPosts = Relationship{
	ID:           uuid.MustParse("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
	LeftModelID:  uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"), // user
	LeftName:     "posts",
	LeftBinding:  HasMany,
	RightModelID: uuid.MustParse("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"), // post
	RightName:    "author",
	RightBinding: BelongsTo,
}

var UserProfile = Relationship{
	ID:           uuid.MustParse("52a31e61-f1d3-4091-8dcf-78236ef84f6f"),
	LeftModelID:  uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"), // user
	LeftName:     "profile",
	LeftBinding:  HasOne,
	RightModelID: uuid.MustParse("66783192-4111-4bd8-95dd-e7da460378df"), // profile
	RightName:    "user",
	RightBinding: BelongsTo,
}

var Profile = Model{
	ID:   uuid.MustParse("66783192-4111-4bd8-95dd-e7da460378df"),
	Name: "profile",
	Attributes: map[string]Attribute{
		"text": Attribute{
			ID:       uuid.MustParse("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		UserProfile,
	},
}

var Post = Model{
	ID:   uuid.MustParse("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	Name: "post",
	Attributes: map[string]Attribute{
		"text": Attribute{
			ID:       uuid.MustParse("b3af6694-b621-43a2-be7f-00956fa505c0"),
			Datatype: String,
		},
	},
	RightRelationships: []Relationship{
		UserPosts,
	},
}

// testing starlark
var Andrew = Datatype{
	ID:          uuid.MustParse("46c0ee11-3943-452d-9420-925dd9be8208"),
	Name:        "andrew",
	FromJSON:    AndrewFromJSON,
	ToJSON:      AndrewToJSON,
	StorageType: StringType,
}

var AndrewFromJSON = Code{
	ID:       uuid.MustParse("aaea187b-d153-4c4a-a7e7-cda599b02ba6"),
	Name:     "andrewFromJson",
	Runtime:  Starlark,
	Function: FromJSON,
	Code: `
def func():
  if args.Value != "Andrew":
  	args.Error = errorf("arg should be Andrew!!!")
func()
`,
}

var AndrewToJSON = Code{
	ID:       uuid.MustParse("dd748a96-b10d-4eff-a582-bf14502d26c4"),
	Name:     "andrewToJson",
	Runtime:  Starlark,
	Function: ToJSON,
	Code: `
def func():
  if args.Value != "Andrew":
  	args.Error = errorf("arg should be Andrew!!!")
func()
`,
}
