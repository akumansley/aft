package db

func AddSampleModels(db DB) {
	tx := db.NewRWTx()
	err := tx.SaveModel(User)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(Profile)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(UserProfile)
	if err != nil {
		panic(err)
	}
	err = tx.SaveModel(Post)
	if err != nil {
		panic(err)
	}
	err = tx.SaveRelationship(UserPosts)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}

var User = Model{
	ID:   MakeModelID("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "user",
	Attributes: []Attribute{
		Attribute{
			Name:     "firstName",
			ID:       MakeID("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			Datatype: String,
		},
		Attribute{
			Name:     "lastName",
			ID:       MakeID("462212e7-dd94-403e-8314-e271fd7ccec9"),
			Datatype: String,
		},
		Attribute{
			Name:     "age",
			ID:       MakeID("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			Datatype: Int,
		},
		Attribute{
			Name:     "emailAddress",
			ID:       MakeID("0fe6bd01-9828-43ac-b004-37620083344d"),
			Datatype: String,
		},
	},
}

var UserPosts = Relationship{
	Name:   "posts",
	ID:     MakeID("28835a3d-6e28-432d-9a9a-b1fe7c468588"),
	Source: User,
	Target: Post,
	Multi:  true,
}

var UserProfile = Relationship{
	Name:   "profile",
	ID:     MakeID("52a31e61-f1d3-4091-8dcf-78236ef84f6f"),
	Source: User,
	Target: Profile,
	Multi:  false,
}

var Profile = Model{
	ID:   MakeModelID("66783192-4111-4bd8-95dd-e7da460378df"),
	Name: "profile",
	Attributes: []Attribute{
		Attribute{
			Name:     "text",
			ID:       MakeID("78fa1725-2b72-4828-8622-f1306b6d0ca7"),
			Datatype: String,
		},
	},
}

var Post = Model{
	ID:   MakeModelID("e25750c8-bb31-41fe-bdec-6bff1dceb2b4"),
	Name: "post",
	Attributes: []Attribute{
		Attribute{
			Name:     "text",
			ID:       MakeID("b3af6694-b621-43a2-be7f-00956fa505c0"),
			Datatype: String,
		},
	},
}
