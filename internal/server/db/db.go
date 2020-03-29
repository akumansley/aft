package db

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/db"
)

func SetupTestData() {
	ObjectTable = &db.TreapTable{}
	ObjectTable.Init()
	for _, obj := range Objects {
		ObjectTable.Put(obj)
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

var ObjectTable db.Table
