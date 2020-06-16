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
	RawData() interface{}
	Map() map[string]interface{}
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	Set(string, interface{}) error
	SetFK(string, uuid.UUID) error
	GetFK(string) (uuid.UUID, error)
	MustGetFK(string) uuid.UUID
	DeepEquals(Record) bool
	DeepCopy() Record
}

// "reflect" based record type
type rRec struct {
	St interface{}
	M  *Model
}

func (r *rRec) RawData() interface{} {
	return r.St
}

func (r *rRec) ID() uuid.UUID {
	id, err := r.Get("id")
	if err != nil {
		panic("Record doesn't have an ID field")
	}
	return id.(uuid.UUID)
}

func (r *rRec) Type() string {
	return r.M.Name
}

func (r *rRec) Model() *Model {
	return r.M
}

func (r *rRec) Get(fieldName string) (interface{}, error) {
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

func (r *rRec) MustGet(fieldName string) interface{} {
	goFieldName := JSONKeyToFieldName(fieldName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)
	if field.IsValid() {
		return field.Interface()
	}
	panic("Key not found")
}

func (r *rRec) Set(name string, value interface{}) error {
	a := r.M.AttributeByName(name)
	d := a.Datatype
	goFieldName := JSONKeyToFieldName(name)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)
	if !field.IsValid() {
		return fmt.Errorf("%w: key %s not found", ErrData, name)
	}
	v, err := d.FromJSON(value)
	if err != nil {
		return err
	}

	if reflect.TypeOf(v) != reflect.TypeOf(storageMap[d.StoredAs]) {
		return fmt.Errorf("%w: Expected type %T and instead found %T", ErrData, v, storageMap[d.StoredAs])
	}
	switch d.StoredAs {
	case BoolStorage:
		b := v.(bool)
		field.SetBool(b)
	case IntStorage:
		i := v.(int64)
		field.SetInt(i)
	case StringStorage:
		s := v.(string)
		field.SetString(s)
	case FloatStorage:
		f := v.(float64)
		field.SetFloat(f)
	case UUIDStorage:
		u := v.(uuid.UUID)
		field.Set(reflect.ValueOf(u))
	}
	return nil
}

func (r *rRec) SetFK(relName string, fkid uuid.UUID) error {
	idFieldName := JSONKeyToRelFieldName(relName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(idFieldName)
	if field.IsValid() {
		v := reflect.ValueOf(fkid)
		field.Set(v)
		return nil
	}
	return fmt.Errorf("%w: key %s not found", ErrData, relName)
}

func (r *rRec) GetFK(relName string) (uuid.UUID, error) {
	idFieldName := JSONKeyToRelFieldName(relName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(idFieldName)
	if field.IsValid() {
		return field.Interface().(uuid.UUID), nil
	}
	return uuid.Nil, fmt.Errorf("%w: key %s not found", ErrData, relName)
}

func (r *rRec) MustGetFK(relName string) uuid.UUID {
	idFieldName := JSONKeyToRelFieldName(relName)
	field := reflect.ValueOf(r.St).Elem().FieldByName(idFieldName)
	if field.IsValid() {
		return field.Interface().(uuid.UUID)
	}
	panic("Key not found")
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
		fieldName := JSONKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageMap[sattr.Datatype.StoredAs]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JSONKeyToFieldName(k)
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageMap[attr.Datatype.StoredAs]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, k, k))}
		fields = append(fields, field)
	}

	for _, b := range m.Bindings() {
		if b.HasField() {
			idFieldName := JSONKeyToRelFieldName(b.Name())
			field := reflect.StructField{
				Name: idFieldName,
				Type: reflect.TypeOf(storageMap[UUID.StoredAs]),
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

func RecordFromParts(st interface{}, m Model) Record {
	return &rRec{St: st, M: &m}
}
