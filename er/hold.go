package er

import (
	"awans.org/aft/er/q"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-immutable-radix"
	"reflect"
)

type Hold struct {
	t *iradix.Tree
}

func New() *Hold {
	return &Hold{t: iradix.New()}
}

func (h *Hold) FindOne(table string, q q.Matcher) (interface{}, error) {
	it := h.t.Root().Iterator()
	it.SeekPrefix([]byte(table))

	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		fmt.Printf("iterating %v\n", val)
		match, err := q.Match(val)
		if err != nil {
			return nil, err
		}
		if match {
			fmt.Printf("matched %v\n", val)
			return val, nil
		}
	}
	fmt.Printf("no match \n")
	return nil, nil
}

func getFieldIf(field string, st interface{}) interface{} {
	k := reflect.ValueOf(st).Kind()
	switch k {
	case reflect.Struct:
		return reflect.ValueOf(st).FieldByName(field).Interface()
	case reflect.Interface:
		return reflect.ValueOf(st).Elem().FieldByName(field).Interface()
	case reflect.Ptr:
		return reflect.ValueOf(st).Elem().FieldByName(field).Interface()

	}
	return nil
}

func getId(st interface{}) uuid.UUID {
	idi := getFieldIf("Id", st)
	id := idi.(uuid.UUID)
	return id
}

func getType(st interface{}) string {
	ti := getFieldIf("Type", st)
	t := ti.(string)
	return t
}

func makeKey(st interface{}) []byte {
	return []byte(fmt.Sprintf("%v/%v", getType(st), getId(st)))
}

func (h *Hold) Insert(object interface{}) {
	fmt.Printf("inserting %v \n", object)
	h.t, _, _ = h.t.Insert(makeKey(object), object)
	h.printTree()
}

func (h *Hold) printTree() {
	fmt.Printf("print tree:\n")
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		fmt.Printf("%v:%v\n", string(k), v)
	}
	fmt.Printf("done printing\n")
}
