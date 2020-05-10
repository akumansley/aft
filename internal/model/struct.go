package model

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
	"strings"
)

// does this really belong in this package
type IncludeResult struct {
	Record         Record
	SingleIncludes map[string]Record
	MultiIncludes  map[string][]Record
}

func (ir IncludeResult) MarshalJSON() ([]byte, error) {
	data := ir.Record.Map()
	for k, v := range ir.SingleIncludes {
		data[k] = v
	}
	for k, v := range ir.MultiIncludes {
		data[k] = v
	}
	return json.Marshal(data)
}

type Record interface {
	Id() uuid.UUID
	Type() string
	Model() *Model
	Map() map[string]interface{}
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
	goFieldName := JsonKeyToFieldName(fieldName)
	return reflect.ValueOf(r.st).Elem().FieldByName(goFieldName).Interface()
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

func (r *rRec) UnmarshalJSON(b []byte) error {
	// just proxy to the inner struct
	if err := json.Unmarshal(b, &r.st); err != nil {
		return err
	}
	return nil
}

func (r *rRec) MarshalJSON() ([]byte, error) {
	// just proxy to the inner struct
	return json.Marshal(r.st)
}

func (r *rRec) Map() map[string]interface{} {
	data := structs.Map(r.st)
	return data
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
		builder.AddField(fieldName, typeMap[sattr.AttrType], fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[attr.AttrType], fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))
	}

	for k, rel := range m.Relationships {
		if rel.HasField() {
			idFieldName := JsonKeyToRelFieldName(k)
			builder.AddField(idFieldName, typeMap[UUID], `json:"-" structs:"-"`)
		}
	}

	b := builder.Build()
	memo[modelName] = b
	st := b.New()

	return &rRec{st: st, m: &m}
}
