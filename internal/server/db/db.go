package db

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/db"
)

func SetupSchema() {
	DB = db.New()
	table := DB.NewTable("objects")
	table.Put(data.ObjectObject)
	table.Put(data.RelationshipObject)
	table.Put(data.AttributeObject)
}

var DB db.DB
