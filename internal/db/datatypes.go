package db

import (
	"awans.org/aft/internal/datatypes"
	"fmt"
	"github.com/google/uuid"
)

type Datatype interface {
	FromJSON(interface{}) (interface{}, error)
	GetID() ID
	Storage() StorageEnumValue
	FillRecord(Record) error
	RecordToDatatype(Record, Tx) (Datatype, error)
}

var datatypeMap map[ID]Datatype = map[ID]Datatype{
	Bool.GetID():              Bool,
	Int.GetID():               Int,
	String.GetID():            String,
	LongText.GetID():          LongText,
	UUID.GetID():              UUID,
	Float.GetID():             Float,
	Runtime.GetID():           Runtime,
	FunctionSignature.GetID(): FunctionSignature,
	StoredAs.GetID():          StoredAs,
}

var storageMap map[StorageEnumValue]interface{} = map[StorageEnumValue]interface{}{
	BoolStorage:   false,
	IntStorage:    int64(0),
	StringStorage: "",
	FloatStorage:  0.0,
	UUIDStorage:   uuid.UUID{},
}

func typeCheck(d Datatype, out interface{}) (interface{}, error) {
	switch d.Storage() {
	case BoolStorage:
		return datatypes.BoolFromJSON(out)
	case IntStorage:
		return datatypes.IntFromJSON(out)
	case StringStorage:
		return datatypes.StringFromJSON(out)
	case FloatStorage:
		return datatypes.FloatFromJSON(out)
	case UUIDStorage:
		return datatypes.UUIDFromJSON(out)
	}
	return nil, fmt.Errorf("Unrecognized storage for datatype")
}

//coreDatatypes
type coreDatatype struct {
	ID        ID
	Name      string
	StoredAs  StorageEnumValue
	Validator Code
}

func (d coreDatatype) FromJSON(arg interface{}) (interface{}, error) {
	c := d.Validator
	out, err := c.Executor.Invoke(c, arg)
	if err != nil {
		panic(err)
	}
	return typeCheck(d, out)
}

func (d coreDatatype) GetID() ID {
	return d.ID
}

func (d coreDatatype) Storage() StorageEnumValue {
	return d.StoredAs
}

func (d coreDatatype) FillRecord(storeDatatype Record) error {
	ew := NewRecordWriter(storeDatatype)
	ew.Set("id", uuid.UUID(d.ID))
	ew.Set("name", d.Name)
	ew.Set("storedAs", uuid.UUID(d.StoredAs.ID))
	ew.Set("enum", false)
	ew.Set("native", true)
	ew.SetFK("validator", d.Validator.ID)
	return ew.err
}

func (d coreDatatype) RecordToDatatype(r Record, tx Tx) (Datatype, error) {
	vk, err := r.GetFK("validator")
	if err != nil {
		return nil, err
	}
	v, err := tx.FindOne(CodeModel.ID, EqID(vk))
	if err != nil {
		return nil, err
	}
	validator, err := RecordToCode(v, tx)
	if err != nil {
		return nil, err
	}
	sa, err := RecordToEnumValue(r, "storedAs", tx)
	if err != nil {
		return nil, err
	}
	ew := NewRecordWriter(r)
	d = coreDatatype{
		ID:        r.ID(),
		Name:      ew.Get("name").(string),
		Validator: validator,
		StoredAs:  StorageEnumValue{sa},
	}
	if ew.err != nil {
		return nil, err
	}
	return d, nil
}

// Created Datatypes
type DatatypeStorage struct {
	ID        ID
	Name      string
	StoredAs  StorageEnumValue
	Validator Code
}

func (d DatatypeStorage) FromJSON(arg interface{}) (interface{}, error) {
	c := d.Validator
	out, err := c.Executor.Invoke(c, arg)
	if err != nil {
		return nil, err
	}
	return typeCheck(d, out)
}

func (d DatatypeStorage) GetID() ID {
	return d.ID
}

func (d DatatypeStorage) Storage() StorageEnumValue {
	return d.StoredAs
}

func (d DatatypeStorage) FillRecord(storeDatatype Record) error {
	ew := NewRecordWriter(storeDatatype)
	ew.Set("id", uuid.UUID(d.ID))
	ew.Set("name", d.Name)
	ew.Set("storedAs", uuid.UUID(d.StoredAs.ID))
	ew.Set("enum", false)
	ew.Set("native", false)
	ew.SetFK("validator", d.Validator.ID)
	return ew.err
}

func (d DatatypeStorage) RecordToDatatype(r Record, tx Tx) (Datatype, error) {
	vk, err := r.GetFK("validator")
	if err != nil {
		return nil, err
	}
	v, err := tx.FindOne(CodeModel.ID, EqID(vk))
	if err != nil {
		return nil, err
	}
	validator, err := RecordToCode(v, tx)
	if err != nil {
		return nil, err
	}
	sa, err := RecordToEnumValue(r, "storedAs", tx)
	if err != nil {
		return nil, err
	}
	ew := NewRecordWriter(r)
	d = DatatypeStorage{
		ID:        r.ID(),
		Name:      ew.Get("name").(string),
		Validator: validator,
		StoredAs:  StorageEnumValue{sa},
	}
	if ew.err != nil {
		return nil, err
	}
	return d, nil
}

//Enums
type Enum struct {
	ID   ID
	Name string
}

func (d Enum) FromJSON(arg interface{}) (interface{}, error) {
	//TODO if it's an enum, verify that it's a UUID.
	// I see two options. One we execute some code here that intelligently
	// checks if this enum id matches the enum values from the table.
	// Second option is that I plum through tx into FromJSON. Second seems worse.
	return datatypes.UUIDFromJSON(arg)
}

func (d Enum) GetID() ID {
	return d.ID
}

func (d Enum) Storage() StorageEnumValue {
	return UUIDStorage
}

func (d Enum) FillRecord(storeDatatype Record) error {
	ew := NewRecordWriter(storeDatatype)
	ew.Set("id", uuid.UUID(d.ID))
	ew.Set("name", d.Name)
	ew.Set("enum", true)
	ew.Set("native", false)
	return ew.err
}

func (d Enum) RecordToDatatype(r Record, tx Tx) (Datatype, error) {
	ew := NewRecordWriter(r)
	d = Enum{
		ID:   r.ID(),
		Name: ew.Get("name").(string),
	}
	if ew.err != nil {
		return nil, ew.err
	}
	return d, nil
}
