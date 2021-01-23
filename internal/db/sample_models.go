package db

func AddSampleModels(db DB) {
	rwtx := db.NewRWTx()
	models := []ModelL{
		User,
		Profile,
		Post,
	}
	for _, m := range models {
		db.AddLiteral(rwtx, m)
	}
	rwtx.Commit()
}

var User = MakeModel(
	MakeID("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	"user",
	[]AttributeL{
		MakeConcreteAttribute(
			MakeID("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			"firstName",
			String,
		),
		MakeConcreteAttribute(
			MakeID("462212e7-dd94-403e-8314-e271fd7ccec9"),
			"lastName",
			String,
		),
		MakeConcreteAttribute(
			MakeID("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			"age",
			Int,
		),
		MakeConcreteAttribute(
			MakeID("0fe6bd01-9828-43ac-b004-37620083344d"),
			"emailAddress",
			String,
		),
	},
	[]RelationshipL{
		UserPosts,
		UserProfile,
	},
	[]ConcreteInterfaceL{},
)

var UserPosts = MakeConcreteRelationship(
	MakeID("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
	"posts",
	true,
	Post,
)

var UserProfile = MakeConcreteRelationship(
	MakeID("52a31e61-f1d3-4091-8dcf-78236ef84f6f"),
	"profile",
	false,
	Profile,
)

var Profile = MakeModel(
	MakeID("66783192-4111-4bd8-95dd-e7da460378df"),
	"profile",
	[]AttributeL{
		MakeConcreteAttribute(
			MakeID("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			"text",
			String,
		),
	},
	[]RelationshipL{},
	[]ConcreteInterfaceL{},
)

func init() {
	Profile.Relationships_ = []RelationshipL{
		ProfileUser,
	}
}

var ProfileUser = MakeReverseRelationship(
	MakeID("ab6510d5-69a8-4240-8f33-22485c3d093e"),
	"user",
	UserProfile,
)

var Post = MakeModel(
	MakeID("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	"post",
	[]AttributeL{
		MakeConcreteAttribute(
			MakeID("b3af6694-b621-43a2-be7f-00956fa505c0"),
			"text",
			String,
		),
	},
	[]RelationshipL{},
	[]ConcreteInterfaceL{},
)
