package db

func AddSampleModels(db DB) {
	models := []Model{
		User,
		Profile,
		Post,
	}
	relationships := []Relationship{
		UserProfile,
		UserPosts,
	}
	tx := db.NewRWTx()
	for _, m := range models {
		err := tx.Schema().SaveModel(m)
		if err != nil {
			panic(err)
		}
	}
	for _, r := range relationships {
		err := tx.Schema().SaveRelationship(r)
		if err != nil {
			panic(err)
		}
	}
	tx.Commit()
}

var User = ModelL{
	ID:   MakeID("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "user",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "firstName",
			ID:       MakeID("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			Datatype: String,
		}.AsAttribute(),
		ConcreteAttributeL{
			Name:     "lastName",
			ID:       MakeID("462212e7-dd94-403e-8314-e271fd7ccec9"),
			Datatype: String,
		}.AsAttribute(),
		ConcreteAttributeL{
			Name:     "age",
			ID:       MakeID("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			Datatype: Int,
		}.AsAttribute(),
		ConcreteAttributeL{
			Name:     "emailAddress",
			ID:       MakeID("0fe6bd01-9828-43ac-b004-37620083344d"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var UserPosts = RelationshipL{
	Name:   "posts",
	ID:     MakeID("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
	Source: User,
	Target: Post,
	Multi:  true,
}.AsRelationship()

var UserProfile = RelationshipL{
	Name:   "profile",
	ID:     MakeID("52a31e61-f1d3-4091-8dcf-78236ef84f6f"),
	Source: User,
	Target: Profile,
	Multi:  false,
}.AsRelationship()

var Profile = ModelL{
	ID:   MakeID("66783192-4111-4bd8-95dd-e7da460378df"),
	Name: "profile",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "text",
			ID:       MakeID("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()

var Post = ModelL{
	ID:   MakeID("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	Name: "post",
	Attributes: []Attribute{
		ConcreteAttributeL{
			Name:     "text",
			ID:       MakeID("b3af6694-b621-43a2-be7f-00956fa505c0"),
			Datatype: String,
		}.AsAttribute(),
	},
}.AsModel()
