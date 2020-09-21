package db

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-immutable-radix"
)

// A hold is an in-memory data structure that stores the indexes and records of a database in aft
//
// It wraps a single big immutable radix tree. Some key indexes include:
//
// - "id/<id>" -> Record
// - "if/<interfaceid>/<id>" -> nil
// - "link/<relid>/<sourceid>/<targetid>" -> nil
// - "rlink/<relid>/<targetid>/<sourceid>" -> nil
//

var (
	ErrNotFound = fmt.Errorf("%w: not found", ErrData)
)

var sep = []byte("/")

type Hold struct {
	t *iradix.Tree
}

type IDIndex struct{}

func (i *IDIndex) makeKey(id ID) []byte {
	var buf bytes.Buffer
	rb, _ := id.Bytes()

	fmt.Fprintf(&buf, "id/%s", rb)
	return buf.Bytes()
}

func PrettyPrintIDIndex(inp []byte) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "id/")
	if bytes.Equal(inp[0:3], buf.Bytes()) {
		return fmt.Sprintf("id/%s", MakeIDFromBytes(inp[3:]))
	}
	return ""
}

func (i *IDIndex) Index(t *iradix.Tree, rec Record) *iradix.Tree {
	k := i.makeKey(rec.ID())
	newTree, _, _ := t.Insert(k, rec)
	return newTree
}

func (i *IDIndex) Get(t *iradix.Tree, id ID) Record {
	k := i.makeKey(id)
	val, found := t.Get(k)
	if !found {
		return nil
	}
	rec := val.(Record)
	return rec
}

func (i *IDIndex) Delete(t *iradix.Tree, rec Record) *iradix.Tree {
	k := i.makeKey(rec.ID())
	t, _, _ = t.Delete(k)
	return t
}

type prefixiter struct {
	it     *iradix.Iterator
	val    interface{}
	err    error
	prefix []byte
}

func newPrefixIterator(t *iradix.Tree, prefix []byte) Iterator {
	it := t.Root().Iterator()
	it.SeekPrefix(prefix)
	return &prefixiter{it, nil, nil, prefix}
}

func (i *prefixiter) Value() interface{} {
	if i.err != nil {
		panic("Called value after err")
	}
	return i.val
}

func (i *prefixiter) Err() error {
	return i.err
}

func (i *prefixiter) Next() bool {
	for k, v, ok := i.it.Next(); ok; k, v, ok = i.it.Next() {
		if !bytes.HasPrefix(k, i.prefix) {
			i.err = Done
			return false
		}

		i.val = v
		return true
	}
	i.err = Done
	return false
}

type IFIndex struct{}

func (i *IFIndex) makePrefix(ifid ID) []byte {
	ifb, _ := ifid.Bytes()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "if/%s", ifb)
	return buf.Bytes()
}

func (i *IFIndex) makeKey(ifid, id ID) []byte {
	rb, _ := id.Bytes()
	ifb, _ := ifid.Bytes()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "if/%s/%s", ifb, rb)
	return buf.Bytes()
}

func PrettyPrintIFIndex(inp []byte) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "if/")
	if bytes.Equal(inp[0:3], buf.Bytes()) {
		return fmt.Sprintf("if/%s/%s", MakeIDFromBytes(inp[3:19]), MakeIDFromBytes(inp[20:36]))
	}
	return ""
}

func (i *IFIndex) makeKeys(rec Record) [][]byte {
	var keys [][]byte
	m := rec.model()
	rid := rec.ID()

	mk := i.makeKey(m.ID(), rid)
	keys = append(keys, mk)

	ifaces, err := m.Implements()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		ik := i.makeKey(iface.ID(), rid)
		keys = append(keys, ik)
	}
	return keys
}

func (i *IFIndex) Index(t *iradix.Tree, rec Record) *iradix.Tree {
	keys := i.makeKeys(rec)
	for _, k := range keys {
		t, _, _ = t.Insert(k, rec)
	}
	return t
}

func (i *IFIndex) Iterator(t *iradix.Tree, ifID ID) Iterator {
	pf := i.makePrefix(ifID)
	return newPrefixIterator(t, pf)
}

func (i *IFIndex) Delete(t *iradix.Tree, rec Record) *iradix.Tree {
	ks := i.makeKeys(rec)
	for _, k := range ks {
		t, _, _ = t.Delete(k)
	}
	return t
}

func newLinkIterator(t *iradix.Tree, rel, id ID, reverse bool) Iterator {
	li := LinkIndex{}
	prefix := li.makePrefix(rel, id, reverse)

	it := t.Root().Iterator()
	it.SeekPrefix(prefix)

	return &linkiter{it, nil, nil, prefix, reverse}
}

// returns IDs pointed to
type linkiter struct {
	it      *iradix.Iterator
	val     interface{}
	err     error
	prefix  []byte
	reverse bool
}

func (i *linkiter) Value() interface{} {
	if i.err != nil {
		panic("Called value after err")
	}
	return i.val
}

func (i *linkiter) Err() error {
	return i.err
}

func (i *linkiter) Next() bool {
	for k, _, ok := i.it.Next(); ok; k, _, ok = i.it.Next() {
		if !bytes.HasPrefix(k, i.prefix) {
			i.err = Done
			return false
		}

		if i.reverse {
			k = k[6+16*2+2 : 6+16*3+2]
		} else {
			k = k[5+16*2+2 : 5+16*3+2]
		}

		id := MakeIDFromBytes(k)
		i.val = id
		return true
	}
	i.err = Done
	return false
}

type LinkIndex struct{}

func (i *LinkIndex) makePrefix(rel, id ID, reverse bool) []byte {
	relb, _ := rel.Bytes()
	idb, _ := id.Bytes()

	prefix := "link"
	if reverse {
		prefix = "rlink"
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s/%s/%s", prefix, relb, idb)
	return buf.Bytes()
}

func PrettyPrintLinkIndex(inp []byte) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "link/")
	if bytes.Equal(inp[0:5], buf.Bytes()) {
		return fmt.Sprintf("link/%s/%s", MakeIDFromBytes(inp[5:21]), MakeIDFromBytes(inp[22:38]))
	}
	return ""
}

func PrettyPrintRlinkIndex(inp []byte) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "rlink/")
	if bytes.Equal(inp[0:6], buf.Bytes()) {
		return fmt.Sprintf("rlink/%s/%s", MakeIDFromBytes(inp[6:22]), MakeIDFromBytes(inp[23:39]))
	}
	return ""
}

func (i *LinkIndex) makeKeys(rel, from, to ID) [][]byte {
	relb, _ := rel.Bytes()
	fromb, _ := from.Bytes()
	tob, _ := to.Bytes()

	var linkBuf, rlinkBuf bytes.Buffer
	fmt.Fprintf(&linkBuf, "link/%s/%s/%s", relb, fromb, tob)
	fmt.Fprintf(&rlinkBuf, "rlink/%s/%s/%s", relb, tob, fromb)

	return [][]byte{
		linkBuf.Bytes(),
		rlinkBuf.Bytes(),
	}
}

func (i *LinkIndex) Link(t *iradix.Tree, from, to, rel ID) *iradix.Tree {
	keys := i.makeKeys(rel, from, to)
	for _, k := range keys {
		t, _, _ = t.Insert(k, nil)
	}

	return t
}

func (i *LinkIndex) Unlink(t *iradix.Tree, from, to, rel ID) *iradix.Tree {
	keys := i.makeKeys(from, to, rel)
	for _, k := range keys {
		t, _, _ = t.Delete(k)
	}

	return t
}

func (i *LinkIndex) Iterator(t *iradix.Tree, rel, id ID, reverse bool) Iterator {
	return newLinkIterator(t, rel, id, reverse)
}

func NewHold() *Hold {
	return &Hold{t: iradix.New()}
}

func (h *Hold) FindOne(modelID ID, q Matcher) (Record, error) {
	ix := IFIndex{}
	it := ix.Iterator(h.t, modelID)

	for it.Next() {
		val := it.Value()
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return nil, err
		}
		if match {
			return rec, nil
		}
	}
	if it.Err() != Done {
		return nil, it.Err()
	}
	return nil, ErrNotFound
}

func (h *Hold) FindMany(modelID ID, q Matcher) ([]Record, error) {
	ix := IFIndex{}
	it := ix.Iterator(h.t, modelID)

	hits := []Record{}
	for it.Next() {
		val := it.Value()
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return hits, err
		}
		if match {
			hits = append(hits, rec)
		}
	}
	if it.Err() != Done {
		return nil, it.Err()
	}
	return hits, nil
}

func (h *Hold) Insert(rec Record) *Hold {
	ifx := IFIndex{}
	t := ifx.Index(h.t, rec)

	idx := IDIndex{}
	t = idx.Index(t, rec)
	return &Hold{t: t}
}

func (h *Hold) Delete(rec Record) *Hold {
	ifx := IFIndex{}
	t := ifx.Delete(h.t, rec)

	idx := IDIndex{}
	t = idx.Delete(t, rec)

	return &Hold{t: t}
}

func (h *Hold) followLinks(id, rel ID, reverse bool) ([]Record, error) {
	lix := LinkIndex{}
	idx := IDIndex{}
	it := lix.Iterator(h.t, rel, id, reverse)

	var hits []Record
	for it.Next() {
		v := it.Value()
		id := v.(ID)
		hit := idx.Get(h.t, id)
		hits = append(hits, hit)
	}
	return hits, nil
}

func (h *Hold) GetLinkedMany(sID, rel ID) ([]Record, error) {
	return h.followLinks(sID, rel, false)
}

func (h *Hold) GetLinkedManyReverse(tID, rel ID) ([]Record, error) {
	return h.followLinks(tID, rel, true)
}

func (h *Hold) followLinksOne(id, rel ID, reverse bool) (Record, error) {
	hits, err := h.followLinks(id, rel, reverse)
	if err != nil {
		return nil, err
	}
	switch len(hits) {
	case 0:
		return nil, ErrNotFound
	case 1:
		return hits[0], nil
	default:
		err = fmt.Errorf("Multi found for to-one rel id: %v rel: %v\n", id, rel)
		panic(err)
	}
}

func (h *Hold) GetLinkedOne(id, rel ID) (Record, error) {
	return h.followLinksOne(id, rel, false)
}

func (h *Hold) GetLinkedOneReverse(id, rel ID) (Record, error) {
	return h.followLinksOne(id, rel, true)
}

type rootiter struct {
	it  *iradix.Iterator
	val interface{}
	err error
}

type item struct {
	k []byte
	v interface{}
}

func (ri *rootiter) Value() interface{} {
	if ri.err != nil {
		panic("Called value after err")
	}
	return ri.val
}

func (ri *rootiter) Err() error {
	return ri.err
}

func (ri *rootiter) Next() bool {
	for k, v, ok := ri.it.Next(); ok; k, v, ok = ri.it.Next() {
		i := item{k, v}
		ri.val = i
		return true
	}
	ri.err = Done
	return false
}

func (h *Hold) Iterator() Iterator {
	it := h.t.Root().Iterator()
	return &rootiter{it: it}
}

func (h *Hold) PrintTree() {
	fmt.Printf("print tree:\n")
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		fmt.Printf("%v:%v\n", string(k), v)
	}
	fmt.Printf("done printing\n")
}

func (h *Hold) String() string {
	var out string
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		var key string
		key = fmt.Sprintf("%s%s%s%s", PrettyPrintIDIndex(k), PrettyPrintIFIndex(k), PrettyPrintLinkIndex(k), PrettyPrintRlinkIndex(k))
		out = fmt.Sprintf("%s%v\n%v\n\n", out, key, v)
	}
	return out
}

func (h *Hold) Link(source, target, rel ID) *Hold {
	li := LinkIndex{}
	t := li.Link(h.t, source, target, rel)
	return &Hold{t: t}
}

func (h *Hold) Unlink(source, target, rel ID) *Hold {
	li := LinkIndex{}
	t := li.Unlink(h.t, source, target, rel)
	return &Hold{t: t}
}
