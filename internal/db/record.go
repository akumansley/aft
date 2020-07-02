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
	ID() ID
	Type() string
	Model() Model
	RawData() interface{}
	Map() map[string]interface{}
	DeepEquals(Record) bool
	DeepCopy() Record

	get(string) (interface{}, error)
	mustGet(string) interface{}
	set(string, interface{}) error
}

// "reflect" based record type
type rRec struct {
	St interface{}
	M  Model
}

func (r *rRec) RawData() interface{} {
	return r.St
}

func (r *rRec) ID() ID {
	id, err := r.get("id")
	if err != nil {
		panic("Record doesn't have an ID field")
	}
	return ID(id.(uuid.UUID))
}

func (r *rRec) Type() string {
	return r.M.Name()
}

func (r *rRec) Model() Model {
	return r.M
}

func (r *rRec) get(fieldName string) (interface{}, error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("bad get: %v on %v %+v - \n", fieldName, r.Type(), r.St)
		}
	}()
	goFieldName := JSONKeyToFieldName(fieldName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)
	if field.IsValid() {
		return field.Interface(), nil
	}
	return nil, fmt.Errorf("%w: key %s not found", ErrData, fieldName)
}

func (r *rRec) mustGet(fieldName string) interface{} {
	goFieldName := JSONKeyToFieldName(fieldName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)
	if field.IsValid() {
		return field.Interface()
	}
	panic("Key not found")
}

func (r *rRec) set(name string, v interface{}) error {
	goFieldName := JSONKeyToFieldName(name)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)

	if !field.IsValid() {
		return fmt.Errorf("%w: key %s not found", ErrData, name)
	}

	field.Set(reflect.ValueOf(v))
	return nil
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

func (r *rRec) DeepCopy() Record {
	newSt := reflect.New(reflect.TypeOf(r.St).Elem())

	val := reflect.ValueOf(r.St).Elem()
	nVal := newSt.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}
	return &rRec{St: newSt.Interface(), M: r.M}
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

var memo = map[string]reflect.Type{}

var storageMap map[ID]interface{} = map[ID]interface{}{
	BoolStorage.ID:   false,
	IntStorage.ID:    int64(0),
	StringStorage.ID: "",
	FloatStorage.ID:  0.0,
	UUIDStorage.ID:   uuid.UUID{},
}

var SystemAttrs = map[string]Attribute{
	"id": ConcreteAttributeL{
		Name:     "id",
		Datatype: UUID,
	}.AsAttribute(),
}

func RecordForModel(m Model) Record {
	modelName := strings.ToLower(m.Name())
	if val, ok := memo[modelName]; ok {
		st := reflect.New(val).Interface()
		return &rRec{St: st, M: m}

	}
	var fields []reflect.StructField

	for k, sattr := range SystemAttrs {
		fieldName := JSONKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageMap[sattr.Storage().ID()]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	// later, maybe we can add validate tags
	attrs, _ := m.Attributes()
	for _, attr := range attrs {
		fieldName := JSONKeyToFieldName(attr.Name())
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageMap[attr.Storage().ID()]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, attr.Name, attr.Name))}
		fields = append(fields, field)
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
	gob.Register(ID{})

	return &rRec{St: st, M: m}
}

func RecordFromParts(st interface{}, m Model) Record {
	return &rRec{St: st, M: m}
}

func JSONKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vID", strings.Title(strings.ToLower(key)))
}

func JSONKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}
