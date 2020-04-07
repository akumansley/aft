package db

import (
	"awans.org/aft/internal/model"
	"github.com/timshannon/bolthold"
)

func InitDB() {
	filename := "db.bolt"
	var err error
	DB, err = bolthold.Open(filename, 0666, nil)
	if err != nil {
		panic(err)
	}
}

func MakeStruct(name string) interface{} {
	if name == "model" {
		return model.Model{}
	} else {
		var m model.Model
		err := DB.FindOne(&m, bolthold.Where("Name").Eq(name))
		if err != nil {
			panic(err)
		}
		return model.StructForModel(m).New()
	}
}

var DB *bolthold.Store
