package db

import (
	"awans.org/aft/internal/data"
)

type DB interface {
	// this needs to be rethought to be oriented around objects
	GetTable(string) Table
	NewTable(string) Table
	Put(string, map[string]interface{}) // can include objects
}

func New() DB {
	db := &mapDB{tables: make(map[string]Table)}

	// bootstrap the db core types
	objects := db.NewTable("objects")
	objects.Put(data.ObjectObject)

	db.NewTable("relationships")
	objects.Put(data.RelationshipObject)

	db.NewTable("attributes")
	objects.Put(data.AttributeObject)
	return db
}

type mapDB struct {
	tables map[string]Table
}

func (db *mapDB) GetTable(name string) Table {
	return db.tables[name]
}

func (db *mapDB) NewTable(name string) Table {
	table := &TreapTable{}
	table.Init()
	db.tables[name] = table
	return table
}

func (db *mapDB) Put(object string, data map[string]interface{}) {
	db.validate(object, data)
	if object == "object" {
		db.putObject(data)
	}
}

func (db *mapDB) validate(object string, data map[string]interface{}) {
	// using the index for this.. just getting something working..
	// results := db.tables["objects"].Query(fmt.Sprintf("name: %v", object))
	// results[0].(map[string]interface{})
	// ..do something with the result

}

func (db *mapDB) putObject(data map[string]interface{}) {

}

func (db *mapDB) Get(object string, id string) interface{} {
	return db.tables[object].Get(id)
}

func (db *mapDB) GetWithRelated(object string, id string, related []string) {

}
