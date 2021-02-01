package db

import (
	"bytes"
	"encoding/binary"
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
	Map() map[string]interface{}
	DeepEquals(Record) bool
	DeepCopy() Record
	String() string
	Version() uint64
	FieldNames() []string

	UnmarshalBinary([]byte) error
	MarshalBinary() ([]byte, error)

	// TODO make these private
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	Set(string, interface{}) error
}

// "reflect" based record type
type rRec struct {
	St interface{}
	I  ID
	s  *Spec
	V  uint64 // spec version
}

func (r *rRec) ID() ID {
	id, err := r.Get("id")
	if err != nil {
		panic("Record doesn't have an ID field")
	}
	return ID(id.(uuid.UUID))
}

func (r *rRec) Version() uint64 {
	return r.V
}

func (r *rRec) Type() string {
	return r.s.InterfaceName
}

func (r *rRec) InterfaceID() ID {
	return r.I
}

func (r *rRec) FieldNames() (result []string) {
	for _, f := range r.s.Fields {
		result = append(result, f.Name)
	}
	return
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
		fmt.Printf("Different types %v %v\n", r.Type(), other.Type())
		return false
	}
	if !reflect.DeepEqual(r.Map(), other.Map()) {
		fmt.Printf("Different values %v %v\n", r.Map(), other.Map())
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

var storageMap map[ID]interface{} = map[ID]interface{}{
	BoolStorage.ID():   false,
	IntStorage.ID():    int64(0),
	StringStorage.ID(): "",
	BytesStorage.ID():  []byte{},
	FloatStorage.ID():  float64(0.0),
	UUIDStorage.ID():   uuid.UUID{},
}

type Spec struct {
	Fields        []Field
	InterfaceID   ID
	InterfaceName string
}

func (s *Spec) StructType() reflect.Type {
	sfs := []reflect.StructField{}
	for _, f := range s.Fields {
		sfs = append(sfs, f.StructField())
	}
	return reflect.StructOf(sfs)
}

type Field struct {
	Name    string
	Storage ID
}

func (f Field) StructField() reflect.StructField {
	fieldName := JSONKeyToFieldName(f.Name)
	sf := reflect.StructField{
		Name: fieldName,
		Type: reflect.TypeOf(storageMap[f.Storage]),
		Tag:  reflect.StructTag(fmt.Sprintf(`json:"%v" structs:"%v"`, f.Name, f.Name))}
	return sf
}

func makeSpec(tx Tx, i Interface) (s *Spec, err error) {
	s = &Spec{
		InterfaceID:   i.ID(),
		InterfaceName: i.Name(),
	}
	attrs, err := i.Attributes(tx)
	if err != nil {
		return
	}

	for _, attr := range attrs {
		sid := attr.Storage(tx).ID()
		if sid == NotStored.ID() {
			continue
		}
		s.Fields = append(s.Fields, Field{Name: attr.Name(), Storage: sid})
	}

	sort.Slice(s.Fields, func(i, j int) bool {
		return s.Fields[i].Name < s.Fields[j].Name
	})
	return
}

func NewBuilder() *Builder {
	return &Builder{
		rtypes:   map[ID][]reflect.Type{},
		registry: map[ID][]*Spec{},
	}
}

type Builder struct {
	rtypes   map[ID][]reflect.Type
	registry map[ID][]*Spec
}

func (b *Builder) InterfaceUpdated(tx Tx, i Interface) error {
	s, err := makeSpec(tx, i)
	if err != nil {
		return err
	}
	b.registerSpec(s, i.ID())
	return nil
}

func (b *Builder) registerSpec(s *Spec, interfaceID ID) {
	b.registry[interfaceID] = append(b.registry[interfaceID], s)
	b.rtypes[interfaceID] = append(b.rtypes[interfaceID], s.StructType())
}

func (b *Builder) getInfo(interfaceID ID) (s *Spec, t reflect.Type, v uint64) {
	if rtypes, ok := b.rtypes[interfaceID]; ok {
		version := len(rtypes) - 1
		return b.registry[interfaceID][version], rtypes[version], uint64(version)
	}
	return
}

func (b *Builder) RecordForInterface(tx Tx, i Interface) (Record, error) {
	spec, sType, v := b.getInfo(i.ID())
	if spec == nil {
		b.InterfaceUpdated(tx, i)
		spec, sType, v = b.getInfo(i.ID())
	}
	st := reflect.New(sType).Interface()
	return &rRec{St: st, I: i.ID(), V: v, s: spec}, nil
}

func (b *Builder) RecordForLiteral(lit Literal) (Record, error) {
	interfaceID := lit.InterfaceID()
	spec, sType, v := b.getInfo(interfaceID)
	if spec == nil {
		s := specFromTaggedLiteral(lit)
		b.registerSpec(s, interfaceID)
		spec, sType, v = b.getInfo(interfaceID)
	}

	st := reflect.New(sType).Interface()
	return &rRec{St: st, I: interfaceID, V: v, s: spec}, nil
}

func (b *Builder) RecordForInterfaceVersion(interfaceID ID, version uint64) (Record, error) {
	if len(b.rtypes[interfaceID]) <= int(version) {
		err := fmt.Errorf("No such interface version: %v %v\n", interfaceID, version)
		panic(err)
	}
	sType := b.rtypes[interfaceID][version]
	spec := b.registry[interfaceID][version]
	st := reflect.New(sType).Interface()
	return &rRec{St: st, I: interfaceID, V: version, s: spec}, nil
}

// rRec must be a record of the correct interface
func (r *rRec) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)

	for _, f := range r.s.Fields {
		if f.Storage == NotStored.ID() {
			continue
		}

		switch f.Storage {
		case BoolStorage.ID():
			var v bool
			err = binary.Read(buf, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.Set(f.Name, v)
		case IntStorage.ID():
			var v int64
			err = binary.Read(buf, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.Set(f.Name, v)
		case StringStorage.ID():
			var byteslen int64
			err = binary.Read(buf, binary.LittleEndian, &byteslen)
			bts := make([]byte, byteslen)
			err = binary.Read(buf, binary.LittleEndian, &bts)
			s := string(bts)
			if err != nil {
				return err
			}
			r.Set(f.Name, s)
		case BytesStorage.ID():
			var byteslen int64
			err = binary.Read(buf, binary.LittleEndian, &byteslen)
			bts := make([]byte, byteslen)
			err = binary.Read(buf, binary.LittleEndian, &bts)
			if err != nil {
				return err
			}
			r.Set(f.Name, bts)
		case FloatStorage.ID():
			var v float64
			err = binary.Read(buf, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			r.Set(f.Name, v)
		case UUIDStorage.ID():
			v := make([]byte, 16)
			err = binary.Read(buf, binary.LittleEndian, &v)
			if err != nil {
				return err
			}
			u, err := uuid.FromBytes(v)
			if err != nil {
				return err
			}
			r.Set(f.Name, u)
		default:
			panic("Invalid storage type")
		}
	}
	return nil
}

func (r *rRec) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, f := range r.s.Fields {
		if f.Storage == NotStored.ID() {
			continue
		}

		val := r.MustGet(f.Name)

		switch f.Storage {
		case BoolStorage.ID():
			err := binary.Write(buf, binary.LittleEndian, val.(bool))
			if err != nil {
				return nil, err
			}
		case IntStorage.ID():
			err := binary.Write(buf, binary.LittleEndian, val.(int64))
			if err != nil {
				return nil, err
			}
		case StringStorage.ID():
			bts := []byte(val.(string))
			btslen := int64(len(bts))
			err := binary.Write(buf, binary.LittleEndian, btslen)
			err = binary.Write(buf, binary.LittleEndian, bts)
			if err != nil {
				return nil, err
			}
		case BytesStorage.ID():
			bts := val.([]byte)
			btslen := int64(len(bts))
			err := binary.Write(buf, binary.LittleEndian, btslen)
			err = binary.Write(buf, binary.LittleEndian, bts)
			if err != nil {
				return nil, err
			}
		case FloatStorage.ID():
			err := binary.Write(buf, binary.LittleEndian, val.(float64))
			if err != nil {
				return nil, err
			}
		case UUIDStorage.ID():
			bytes, _ := val.(uuid.UUID).MarshalBinary()
			err := binary.Write(buf, binary.LittleEndian, bytes)
			if err != nil {
				return nil, err
			}
		default:
			panic("Invalid storage type")
		}
	}
	return buf.Bytes(), nil
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
