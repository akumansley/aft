package db

import (
	"bytes"
	"fmt"
)

var ErrNotFound = fmt.Errorf("%w: not found", ErrData)

type KVIterator interface {
	Next() bool
	Err() error
	Entry() ([]byte, interface{})
	Key() []byte
	Value() interface{}
}

type IDIterator interface {
	Next() bool
	Err() error
	ID() ID
}

type LinkIterator interface {
	Next() bool
	Err() error
	TargetID() ID
	SourceID() ID
	RelID() ID
}

// migration methods

func (h *hold) dropImplements(modelID, interfaceID ID) *hold {
	h.implements.IndexDelete(modelID, interfaceID)
	ifaceIter := h.iface.Iterator(modelID)
	for ifaceIter.Next() {
		recID := ifaceIter.ID()
		h.iface.IndexDelete(interfaceID, recID)
	}
	return h.snap()
}

func (h *hold) addImplements(modelID, interfaceID ID) *hold {
	h.implements.Index(modelID, interfaceID)
	ifaceIter := h.iface.Iterator(modelID)
	for ifaceIter.Next() {
		recID := ifaceIter.ID()
		h.iface.Index(interfaceID, recID)
	}
	return h.snap()
}

func (h *hold) dropRel(sourceInterfaceID, targetInterfaceID, relID ID) (*hold, error) {
	h.rel.IndexDelete(sourceInterfaceID, relID)
	h.rrel.IndexDelete(targetInterfaceID, relID)
	linkIter := h.link.PrefixIterator(relID)
	for linkIter.Next() {
		targetID := linkIter.TargetID()
		sourceID := linkIter.SourceID()
		h.link.IndexDelete(relID, sourceID, targetID)
		h.rlink.IndexDelete(relID, targetID, sourceID)
	}
	if linkIter.Err() != Done {
		return nil, linkIter.Err()
	}

	rlinkIter := h.rlink.PrefixIterator(relID)
	for rlinkIter.Next() {
		targetID := rlinkIter.TargetID()
		sourceID := rlinkIter.SourceID()
		h.rlink.IndexDelete(relID, sourceID, targetID)
		h.link.IndexDelete(relID, targetID, sourceID)
	}
	if rlinkIter.Err() != Done {
		return nil, rlinkIter.Err()
	}
	return h.snap(), nil
}

func (h *hold) addRel(sourceInterfaceID, targetInterfaceID, relID ID) *hold {
	h.rel.Index(sourceInterfaceID, relID)
	h.rrel.Index(targetInterfaceID, relID)
	return h.snap()
}

// end of migration methods

// read implements and write from record -> iface
func (h *hold) indexInterfaces(rec Record) error {
	modelID := rec.InterfaceID()
	recID := rec.ID()
	iter := h.implements.Iterator(modelID)
	for iter.Next() {
		implementsID := iter.ID()
		h.iface.Index(implementsID, recID)
	}
	if iter.Err() != Done {
		return iter.Err()
	}
	h.iface.Index(modelID, recID)
	return nil
}

// read implements and drop from record -> iface
// read rel and drop from link/rlink
func (h *hold) cascadeDelete(rec Record) error {
	// clean up interface index entries
	modelID := rec.InterfaceID()
	recID := rec.ID()
	implementsIter := h.implements.Iterator(modelID)
	for implementsIter.Next() {
		implementsID := implementsIter.ID()
		h.iface.IndexDelete(implementsID, recID)
	}
	if implementsIter.Err() != Done {
		return implementsIter.Err()
	}
	h.iface.IndexDelete(modelID, recID)

	// clean up rels
	relIter := h.rel.Iterator(modelID)
	for relIter.Next() {
		relID := relIter.ID()
		linkIter := h.link.Iterator(relID, recID)
		for linkIter.Next() {
			targetID := linkIter.TargetID()
			h.link.IndexDelete(relID, recID, targetID)
		}
		if linkIter.Err() != Done {
			return linkIter.Err()
		}
	}
	rrelIter := h.rrel.Iterator(modelID)
	for rrelIter.Next() {
		rrelID := rrelIter.ID()
		rlinkIter := h.rlink.Iterator(rrelID, recID)
		for rlinkIter.Next() {
			targetID := rlinkIter.TargetID()
			h.rlink.IndexDelete(rrelID, recID, targetID)
		}
		if rlinkIter.Err() != Done {
			return rlinkIter.Err()
		}
	}
	return nil
}

type hold struct {
	// mutable kv store
	// writes go here
	kv kv

	// immutable snapshot of kv
	// reads come from here
	s snapshot

	// schema indexes
	// modelID -> relID
	rel *index2 // 'e'
	// modelID -> relID, but the reverse ones
	rrel *index2 // 's'

	// modelID -> interfaceID
	implements *index2 // 'm'

	// records
	// id -> record
	records *recordIndex // 'r'
	// ifaceID -> recordID
	iface *index2 // 'f'

	// link
	rlink *linkIndex // 'v'
	link  *linkIndex // 'l'

}

func NewHold() *hold {
	kv := newIRadixKV()
	s := kv.Snapshot()
	return &hold{
		kv:         kv,
		s:          s,
		rel:        newIndex2(newView(kv, s, 'e')),
		rrel:       newIndex2(newView(kv, s, 's')),
		implements: newIndex2(newView(kv, s, 'm')),
		records:    newRecordIndex(newView(kv, s, 'r')),
		iface:      newIndex2(newView(kv, s, 'f')),
		rlink:      newLinkIndex(newView(kv, s, 'v')),
		link:       newLinkIndex(newView(kv, s, 'l')),
	}
}

// snap returns a copy of a hold, updating the snapshot
// TODO this is a lot of allocations
func (h *hold) snap() *hold {
	kv := h.kv
	s := h.kv.Snapshot()

	return &hold{
		kv:         kv,
		s:          s,
		rel:        newIndex2(newView(kv, s, 'e')),
		rrel:       newIndex2(newView(kv, s, 's')),
		implements: newIndex2(newView(kv, s, 'm')),
		records:    newRecordIndex(newView(kv, s, 'r')),
		iface:      newIndex2(newView(kv, s, 'f')),
		rlink:      newLinkIndex(newView(kv, s, 'v')),
		link:       newLinkIndex(newView(kv, s, 'l')),
	}
}

func (h *hold) FindOne(interfaceID ID, q Matcher) (result Record, err error) {
	iter := h.iface.Iterator(interfaceID)
	for iter.Next() {
		id := iter.ID()
		result, err = h.records.Get(id)
		if err != nil {
			return
		}
		var match bool
		match, err = q.Match(result)
		if err != nil {
			return nil, err
		}
		if match {
			return
		}
	}
	if err = iter.Err(); err != nil {
		return
	}
	return nil, ErrNotFound
}

func (h *hold) FindMany(interfaceID ID, q Matcher) (results []Record, err error) {
	iter := h.iface.Iterator(interfaceID)
	for iter.Next() {
		id := iter.ID()
		var rec Record
		rec, err = h.records.Get(id)
		if err != nil {
			return
		}
		var match bool
		match, err = q.Match(rec)
		if err != nil {
			return nil, err
		}
		if match {
			results = append(results, rec)
		}
	}
	if iter.Err() != Done {
		return nil, iter.Err()
	}
	return
}

func (h *hold) Insert(rec Record) *hold {
	h.records.Index(rec)
	h.indexInterfaces(rec)
	return h.snap()
}

func (h *hold) Delete(rec Record) *hold {
	h.records.IndexDelete(rec)
	h.cascadeDelete(rec)
	return h.snap()
}

func (h *hold) GetLinkedMany(sourceID, rel ID) (results []Record, err error) {
	iter := h.link.Iterator(rel, sourceID)
	for iter.Next() {
		id := iter.TargetID()
		var rec Record
		rec, err = h.records.Get(id)
		if err != nil {
			return
		}
		results = append(results, rec)
	}
	if iter.Err() != Done {
		// we could just return results, err
		// this is more.. defensive
		return nil, iter.Err()
	}
	return
}

func (h *hold) GetLinkedOne(sourceID, rel ID) (result Record, err error) {
	iter := h.link.Iterator(rel, sourceID)
	for iter.Next() {
		id := iter.TargetID()
		result, err = h.records.Get(id)
		return
	}
	if iter.Err() != Done {
		return nil, iter.Err()
	}
	return nil, ErrNotFound
}

func (h *hold) GetLinkedManyReverse(targetID, rel ID) (results []Record, err error) {
	iter := h.rlink.Iterator(rel, targetID)
	for iter.Next() {
		id := iter.TargetID()
		var rec Record
		rec, err = h.records.Get(id)
		if err != nil {
			return
		}
		results = append(results, rec)
	}
	if iter.Err() != Done {
		return nil, iter.Err()
	}
	return
}

func (h *hold) GetLinkedOneReverse(targetID, rel ID) (result Record, err error) {
	iter := h.rlink.Iterator(rel, targetID)
	for iter.Next() {
		id := iter.TargetID()
		result, err = h.records.Get(id)
		return
	}
	if iter.Err() != Done {
		return nil, iter.Err()
	}
	return nil, ErrNotFound
}

func (h *hold) Link(sourceID, targetID, relID ID) *hold {
	h.link.Index(relID, sourceID, targetID)
	h.rlink.Index(relID, targetID, sourceID)
	return h.snap()
}
func (h *hold) Unlink(sourceID, targetID, relID ID) *hold {
	h.link.IndexDelete(relID, sourceID, targetID)
	h.rlink.IndexDelete(relID, targetID, sourceID)
	return h.snap()
}

func (h *hold) Iterator() KVIterator {
	return h.kv.PrefixIterator([]byte{})
}

func (h *hold) String() string {
	var buf bytes.Buffer
	buf.WriteString("hold{")
	iter := h.Iterator()
	for iter.Next() {
		buf.Write(iter.Key())
		buf.WriteByte('\n')
	}
	buf.WriteString("}")
	return buf.String()
}

// views are a namespacing convenience for index management

func newView(kv kv, s snapshot, prefix byte) *view {
	return &view{
		kv:     kv,
		s:      s,
		prefix: prefix,
	}
}

type view struct {
	kv     kv
	s      snapshot
	prefix byte
}

func (v view) Put(key []byte, value interface{}) error {
	key[0] = v.prefix
	return v.kv.Put(key, value)
}

func (v view) Get(key []byte) (interface{}, error) {
	key[0] = v.prefix
	return v.s.Get(key)
}

func (v view) Delete(key []byte) error {
	key[0] = v.prefix
	return v.kv.Delete(key)
}

func (v view) PrefixIterator(key []byte) KVIterator {
	key[0] = v.prefix
	return v.s.PrefixIterator(key)
}
