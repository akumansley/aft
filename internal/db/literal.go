package db

import (
	"reflect"
)

func MarshalRecord(v interface{}, m Model) (rec Record, err error) {
	rec = RecordForModel(m)

	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	attrMap := map[string]Attribute{}
	for _, a := range attrs {
		attrMap[a.Name()] = a
	}

	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v).Elem()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		recFieldName := field.Tag.Get("record")
		attr, ok := attrMap[recFieldName]
		if !ok {
			panic("failed to marshal struct to record")
		}
		vIf := vVal.Field(i).Interface()
		var f Function
		f, err = attr.Datatype().FromJSON()
		if err != nil {
			return
		}
		var converted interface{}
		converted, err = f.Call(vIf)
		if err != nil {
			return
		}
		rec.Set(recFieldName, converted)
	}
	return rec, nil
}
