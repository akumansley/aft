package db

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
)

func (i Include) Resolve(db DB, st interface{}) interface{} {
	for _, inc := range i.Includes {
		db.Resolve(st, inc)
	}
	return st
}

func getFK(st interface{}, key string) uuid.UUID {
	fieldName := model.JsonKeyToRelFieldName(key)
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField(fieldName).Interface().(uuid.UUID)
	return id
}

func (db holdDB) Resolve(st interface{}, i Inclusion) {
	id := getId(st)
	var m q.Matcher
	rel := i.Relationship
	backRel := getBackref(db, rel)
	var related interface{}
	switch rel.RelType {
	case model.HasOne:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := db.h.IterMatches(rel.TargetModel, m)
		var hits []interface{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on hasOne")
		}
		related = hits[0]
	case model.BelongsTo:
		// FK on this side
		thisFK := getFK(st, backRel.TargetRel)
		m = q.Eq("Id", thisFK)
		mi := db.h.IterMatches(rel.TargetModel, m)
		var hits []interface{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on belongTO")
		}
		related = hits[0]
	case model.HasMany:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := db.h.IterMatches(rel.TargetModel, m)
		hits := []interface{}{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		related = hits
	case model.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	setRelated(st, backRel.TargetRel, related)

}

func setRelated(st interface{}, key string, val interface{}) {
	fieldName := model.JsonKeyToFieldName(key)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	v := reflect.ValueOf(&val)
	field.Set(v)
}
