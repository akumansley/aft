package db

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

type Record interface {
	ID() uuid.UUID
	Type() string
	Model() *Model
	Map() map[string]interface{}
	Get(string) interface{}
	Set(string, interface{}) error
	SetFK(string, uuid.UUID)
	GetFK(string) uuid.UUID
	DeepEquals(Record) bool
}

// "reflect" based record type
type rRec struct {
	St interface{}
	M  *Model
}

func (r *rRec) ID() uuid.UUID {
	return r.Get("id").(uuid.UUID)
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

func (r *rRec) Set(name string, value interface{}) error {
	a := r.M.AttributeByName(name)
	goFieldName := JsonKeyToFieldName(name)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)
	parsedValue, err := CallFunc(a.Datatype.FromJson, value)
	if err != nil {
		return err
	}

	switch parsedValue.(type) {
	case bool:
		b := parsedValue.(bool)
		field.SetBool(b)
	case int64:
		i := parsedValue.(int64)
		field.SetInt(i)
	case string:
		s := parsedValue.(string)
		field.SetString(s)
	case float64:
		f := parsedValue.(float64)
		field.SetFloat(f)
	case uuid.UUID:
		u := parsedValue.(uuid.UUID)
		v := reflect.ValueOf(u)
		field.Set(v)
	}
	return nil
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
	//iterate over the key and call toJson on each one.
	// just proxy to the inner struct
	// TODO
	return json.Marshal(r.St)
}

func (r *rRec) Map() map[string]interface{} {
	data := structs.Map(r.St)
	return data
}

var memo = map[string]reflect.Type{}

var SystemAttrs = map[string]Attribute{
	"id": Attribute{
		Datatype: UUID,
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
			Type: reflect.TypeOf(storageTypeMap[sattr.Datatype.Type]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageTypeMap[attr.Datatype.Type]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	for _, b := range m.Bindings() {
		if b.HasField() {
			idFieldName := JsonKeyToRelFieldName(b.Name())
			field := reflect.StructField{
				Name: idFieldName,
				Type: reflect.TypeOf(storageTypeMap[UUID.Type]),
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
