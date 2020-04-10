package db

import (
	"awans.org/aft/internal/model"
	"github.com/timshannon/bolthold"
	"io/ioutil"
	"strings"
)

func New() DB {
	tmpfile, err := ioutil.TempFile("", "db.bolt")
	if err != nil {
		panic(err)
	}

	db, err := bolthold.Open(tmpfile.Name(), 0666, nil)
	if err != nil {
		panic(err)
	}
	return DB{db: db}
}

type DB struct {
	db *bolthold.Store
}

func (db DB) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	if modelName == "model" {
		return model.Model{}
	} else {
		var m model.Model
		err := db.db.FindOne(&m, bolthold.Where("Name").Eq(modelName))
		if err != nil {
			panic(err)
		}
		return model.StructForModel(m).New()
	}
}

func (db DB) GetModel(modelName string) model.Model {
	modelName = strings.ToLower(modelName)
	var m model.Model
	err := db.db.FindOne(&m, bolthold.Where("Name").Eq(modelName))
	if err != nil {
		panic(err)
	}
	return m
}
