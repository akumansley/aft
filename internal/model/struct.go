package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
	"strings"
)

type Record interface {
	Id() uuid.UUID
	Type() string
	Model() *Model
	Get(string) interface{}
	Set(string, interface{})
	SetFK(string, uuid.UUID)
	GetFK(string) uuid.UUID
}

// "reflect" (dynamicstruct) based record type
type rRec struct {
	st interface{}
	m  *Model
}

func (r *rRec) Id() uuid.UUID {
	return r.Get("Id").(uuid.UUID)
}

func (r *rRec) Type() string {
	return r.m.Name
}

func (r *rRec) Model() *Model {
	return r.m
}

func (r *rRec) Get(fieldName string) interface{} {
	return reflect.ValueOf(r.st).Elem().FieldByName(fieldName).Interface()
}

func (r *rRec) Set(fieldName string, value interface{}) {
	a, ok := r.m.Attributes[fieldName]
	if !ok {
		a, ok = SystemAttrs[fieldName]
	}
	// maybe refactor SetField to be inside here
	a.SetField(fieldName, value, r.st)
}
func (r *rRec) SetFK(relName string, fkid uuid.UUID) {
	idFieldName := JsonKeyToRelFieldName(relName)
	field := reflect.ValueOf(r.st).Elem().FieldByName(idFieldName)
	v := reflect.ValueOf(fkid)
	field.Set(v)
}

func (r *rRec) GetFK(relName string) uuid.UUID {
	idFieldName := JsonKeyToRelFieldName(relName)
	idif := reflect.ValueOf(r.st).Elem().FieldByName(idFieldName).Interface()
	u := idif.(uuid.UUID)
	return u
}

var typeMap map[AttrType]interface{} = map[AttrType]interface{}{
	Int:    int64(0),
	String: "",
	Text:   "",
	Float:  0.0,
	Enum:   int64(0),
	UUID:   uuid.UUID{},
}

var memo = map[string]dynamicstruct.DynamicStruct{}

var SystemAttrs = map[string]Attribute{
	"id": Attribute{
		AttrType: UUID,
	},
	"type": Attribute{
		AttrType: String,
	},
}

func RecordForModel(m Model) Record {
	modelName := strings.ToLower(m.Name)
	if val, ok := memo[modelName]; ok {
		st := val.New()
		return &rRec{st: st, m: &m}

	}

	builder := dynamicstruct.NewStruct()

	for k, sattr := range SystemAttrs {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[sattr.AttrType], fmt.Sprintf(`json:"%v"`, k))
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[attr.AttrType], fmt.Sprintf(`json:"%v"`, k))
	}

	for k, rel := range m.Relationships {
		if rel.HasField() {
			idFieldName := JsonKeyToRelFieldName(k)
			builder.AddField(idFieldName, typeMap[UUID], `json:"-"`)
		}
		colFieldName := JsonKeyToFieldName(k)
		var i interface{}
		builder.AddField(colFieldName, &i, fmt.Sprintf(`json:"%v,omitempty"`, k))
	}

	b := builder.Build()
	memo[modelName] = b
	st := b.New()

	return &rRec{st: st, m: &m}
}
