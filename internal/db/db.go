package db

import ()

type DB interface {
	GetTable(string) Table
	NewTable(string) Table
}

func New() DB {
	return &mapDB{tables: make(map[string]Table)}
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
