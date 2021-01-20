package db

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/structs"
	"github.com/google/uuid"
)

type Record interface {
	ID() ID
	Type() string
	InterfaceID() ID
	Interface() Interface
	RawData() interface{}
	Map() map[string]interface{}
	DeepEquals(Record) bool
	DeepCopy() Record
	String() string

	model() Model

	// TODO make these private
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	Set(string, interface{}) error
}

// "reflect" based record type
type rRec struct {
	St interface{}
	I  Interface
}

func (r *rRec) RawData() interface{} {
	return r.St
}

func (r *rRec) ID() ID {
	id, err := r.Get("id")
	if err != nil {
		panic("Record doesn't have an ID field")
	}
	return ID(id.(uuid.UUID))
}

func (r *rRec) Type() string {
	return r.I.Name()
}

func (r *rRec) InterfaceID() ID {
	return r.I.ID()
}

func (r *rRec) Interface() Interface {
	return r.I
}

func (r *rRec) model() Model {
	m, ok := r.I.(Model)
	if ok {
		return m
	}
	err := fmt.Errorf("model() on non-concrete record %v", r.Map())
	panic(err)
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
	err := fmt.Errorf("key %v not found on %v", fieldName, r.Map())
	panic(err)
}

func (r *rRec) Set(name string, v interface{}) error {
	goFieldName := JSONKeyToFieldName(name)
	field := reflect.ValueOf(r.St).Elem().FieldByName(goFieldName)

	if !field.IsValid() {
		err := fmt.Errorf("%w: key %s not found", ErrData, name)
		panic(err)
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
	return &rRec{St: newSt.Interface(), I: r.I}
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

func (r *rRec) String() string {
	bytes, _ := r.MarshalJSON()
	return string(bytes)
}

func (r *rRec) Map() map[string]interface{} {
	data := structs.Map(r.St)
	data["type"] = r.Type()
	return data
}

var memo = map[ID]reflect.Type{}

var storageMap map[ID]interface{} = map[ID]interface{}{
	BoolStorage.ID():   false,
	IntStorage.ID():    int64(0),
	StringStorage.ID(): "",
	BytesStorage.ID():  []byte{},
	FloatStorage.ID():  0.0,
	UUIDStorage.ID():   uuid.UUID{},
}

func newID(rec Record) error {
	u := uuid.New()
	err := rec.Set("id", u)
	return err
}

func NewRecord(m Interface) Record {
	rec := RecordForModel(m)
	newID(rec)
	return rec
}

func interfaceUpdated(iface Interface) {
	delete(memo, iface.ID())
	RecordForModel(iface)
}

func RecordForModel(i Interface) Record {
	ifaceID := i.ID()
	if val, ok := memo[ifaceID]; ok {
		st := reflect.New(val).Interface()
		return &rRec{St: st, I: i}

	}
	var fields []reflect.StructField

	// later, maybe we can add validate tags
	attrs, err := i.Attributes()
	if err != nil {
		panic(err)
	}
	for _, attr := range attrs {
		if attr.Storage().ID() == NotStored.ID() {
			continue
		}
		fieldName := JSONKeyToFieldName(attr.Name())
		field := reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(storageMap[attr.Storage().ID()]),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, attr.Name(), attr.Name()))}
		fields = append(fields, field)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	sType := reflect.StructOf(fields)
	memo[ifaceID] = sType
	st := reflect.New(sType).Interface()

	// can't see a way around this
	// for now -- it's a hack for goblog.go
	// to be able to gob encode / decode
	// these generated types
	gob.Register(st)
	gob.Register(&st)
	gob.Register(&rRec{})
	gob.Register(ID{})

	return &rRec{St: st, I: i}
}

func RecordFromParts(st interface{}, i Interface) Record {
	return &rRec{St: st, I: i}
}

func JSONKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vID", strings.Title(strings.ToLower(key)))
}

func JSONKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

var StoredAs = MakeEnum(
	MakeID("30a04b8c-720a-468e-8bc6-6ff101e412b3"),
	"storedAs",
	[]EnumValueL{
		BoolStorage,
		IntStorage,
		StringStorage,
		BytesStorage,
		FloatStorage,
		UUIDStorage,
		NotStored,
	})

var NotStored = MakeEnumValue(
	MakeID("e0f86fe9-10ea-430b-a393-b01957a3eabf"),
	"notStored",
)

var BoolStorage = MakeEnumValue(
	MakeID("4f71b3af-aad5-422a-8729-e4c0273aa9bd"),
	"bool",
)

var IntStorage = MakeEnumValue(
	MakeID("14b3d69a-a940-4418-aca1-cec12780b449"),
	"int",
)

var StringStorage = MakeEnumValue(
	MakeID("200630e4-6724-406e-8218-6161bcefb3d4"),
	"string",
)

var BytesStorage = MakeEnumValue(
	MakeID("bc7a618f-e87a-4044-a451-9e239212fe2e"),
	"bytes",
)

var FloatStorage = MakeEnumValue(
	MakeID("ef9995c7-2881-44de-98ff-8960df0e5046"),
	"float",
)

var UUIDStorage = MakeEnumValue(
	MakeID("4d744a2c-e3f3-4a8b-b645-0af46b0235ae"),
	"uuid",
)
