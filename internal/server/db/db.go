package db

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/db"
)

func SetupTestData() {
	DB = db.New()
	table := DB.NewTable("objects")
	for _, obj := range Objects {
		table.Put(obj)
	}
}

var Objects = []data.Object{
	data.Object{
		Id:   "Cekw67uyMpBGZLRP2HFVbe",
		Name: "Test",
		Fields: []data.Field{
			data.Field{
				Name: "f1",
				Type: data.Text,
			},
		},
	},
	data.Object{
		Id:   "6R7VqaQHbzC1xwA5UueGe6",
		Name: "Cool",
		Fields: []data.Field{
			data.Field{
				Name: "f5",
				Type: data.Int,
			},
		},
	},
}

var DB db.DB
