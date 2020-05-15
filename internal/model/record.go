package model

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"reflect"
	"sort"
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
	DeepEquals(Record) bool
}

// "reflect" based record type
type rRec struct {
	St interface{}
	M  *Model
}

func (r *rRec) Id() uuid.UUID {
	return r.Get("Id").(uuid.UUID)
}

func (r *rRec) Type() string {
	return r.M.Name
}

func (r *rRec) Model() *Model {
	return r.M
}

func (r *rRec) Get(fieldName string) interface{} {
	goFieldName := JsonKeyToFieldName(fieldName)
	return reflect.ValueOf(r.St).Elem().FieldByName(goFieldName).Interface()
}

func (r *rRec) Set(fieldName string, value interface{}) {
	a, ok := r.M.Attributes[fieldName]
	if !ok {
		a, ok = SystemAttrs[fieldName]
	}
	// maybe refactor SetField to be inside here
	a.SetField(fieldName, value, r.St)
}
func (r *rRec) SetFK(relName string, fkid uuid.UUID) {
	idFieldName := JsonKeyToRelFieldName(relName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(idFieldName)
	v := reflect.ValueOf(fkid)
	field.Set(v)
}

func (r *rRec) GetFK(relName string) uuid.UUID {
	idFieldName := JsonKeyToRelFieldName(relName)
	idif := reflect.ValueOf(r.St).Elem().FieldByName(idFieldName).Interface()
	u := idif.(uuid.UUID)
	return u
}

func (r *rRec) DeepEquals(other Record) bool {
	if r.Type() != other.Type() {
		return false
	}
	if !reflect.DeepEqual(r.Map(), other.Map()) {
		return false
	}
	return true
}

func (r *rRec) UnmarshalJSON(b []byte) error {
	// just proxy to the inner struct
	if err := json.Unmarshal(b, &r.St); err != nil {
		return err
	}
	return nil
}

func (r *rRec) MarshalJSON() ([]byte, error) {
	// just proxy to the inner struct
	return json.Marshal(r.St)
}

func (r *rRec) Map() map[string]interface{} {
	data := structs.Map(r.St)
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

var memo = map[string]reflect.Type{}

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
		st := reflect.New(val).Interface()
		return &rRec{St: st, M: &m}

	}
	var fields []reflect.StructField

	for k, sattr := range SystemAttrs {
		fieldName := JsonKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(typeMap[sattr.AttrType]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(typeMap[attr.AttrType]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	for k, rel := range m.Relationships {
		if rel.HasField() {
			idFieldName := JsonKeyToRelFieldName(k)
			field := reflect.StructField{
				Name: idFieldName,
				Type: reflect.TypeOf(typeMap[UUID]),
				Tag:  reflect.StructTag(`json:"-" structs:"-"`)}
			fields = append(fields, field)
		}
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	sType := reflect.StructOf(fields)

	memo[modelName] = sType
	st := reflect.New(sType).Interface()

	// can't see a way around this
	// for now -- it's a hack for goblog.go
	// to be able to gob encode / decode
	// these generated types
	gob.Register(st)
	gob.Register(&st)
	gob.Register(&rRec{})

	return &rRec{St: st, M: &m}
}
