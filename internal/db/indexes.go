package db

import "bytes"

const placeholder byte = 0

func makeKey(id1 ID) []byte {
	var key bytes.Buffer
	key.WriteByte(placeholder)
	key.Write(id1.Bytes())
	return key.Bytes()
}

func join2(id1, id2 ID) []byte {
	var key bytes.Buffer
	key.WriteByte(placeholder)
	key.Write(id1.Bytes())
	key.Write(id2.Bytes())
	return key.Bytes()
}

func split2(bytes []byte) (id1, id2 ID) {
	id1Bytes := bytes[1:17]
	id1 = MakeIDFromBytes(id1Bytes)
	id2Bytes := bytes[17:33]
	id2 = MakeIDFromBytes(id2Bytes)
	return id1, id2
}

func join3(id1, id2, id3 ID) []byte {
	var key bytes.Buffer
	key.WriteByte(placeholder)
	key.Write(id1.Bytes())
	key.Write(id2.Bytes())
	key.Write(id3.Bytes())
	return key.Bytes()
}

func split3(bytes []byte) (id1, id2, id3 ID) {
	id1Bytes := bytes[1:17]
	id1 = MakeIDFromBytes(id1Bytes)
	id2Bytes := bytes[17:33]
	id2 = MakeIDFromBytes(id2Bytes)
	id3Bytes := bytes[33:49]
	id3 = MakeIDFromBytes(id3Bytes)
	return id1, id2, id3
}

// the coordination to lookup the right
// interfaces is handled in the librarian
// i've tried hard to keep the indexes "dumb"

func newIndex2(v *view) *index2 {
	return &index2{
		view: v,
	}
}

type index2 struct {
	view *view
}

func (i *index2) Index(id1, id2 ID) error {
	key := join2(id1, id2)
	return i.view.Put(key, nil)
}

func (i *index2) IndexDelete(id1, id2 ID) error {
	key := join2(id1, id2)
	return i.view.Delete(key)
}

func (i *index2) Iterator(id1 ID) IDIterator {
	prefix := makeKey(id1)
	pi := i.view.PrefixIterator(prefix)
	return &iditer{pi}
}

type iditer struct {
	KVIterator
}

func (i *iditer) ID() ID {
	k := i.Key()
	_, id2 := split2(k)
	return id2
}

func newRecordIndex(v *view) *recordIndex {
	return &recordIndex{
		view: v,
	}
}

type recordIndex struct {
	view *view
}

func (r *recordIndex) Index(rec Record) error {
	return r.view.Put(makeKey(rec.ID()), rec)
}

func (r *recordIndex) IndexDelete(rec Record) error {
	return r.view.Delete(makeKey(rec.ID()))
}

func (r *recordIndex) Get(id ID) (Record, error) {
	val, err := r.view.Get(makeKey(id))
	if err != nil {
		return nil, err
	}
	rec, ok := val.(Record)
	if !ok {
		panic("Non-record value in recordIndex")
	}
	return rec, nil
}

func newLinkIndex(v *view) *linkIndex {
	return &linkIndex{
		view: v,
	}
}

type linkIndex struct {
	view *view
}

func (l *linkIndex) Index(relID, sourceID, targetID ID) error {
	key := join3(relID, sourceID, targetID)
	return l.view.Put(key, nil)
}

func (l *linkIndex) IndexDelete(relID, sourceID, targetID ID) error {
	key := join3(relID, sourceID, targetID)
	return l.view.Delete(key)
}

func (l *linkIndex) PrefixIterator(relID ID) LinkIterator {
	prefix := makeKey(relID)
	pi := l.view.PrefixIterator(prefix)
	return &iditer3{pi}
}

func (l *linkIndex) Iterator(relID, sourceID ID) LinkIterator {
	prefix := join2(relID, sourceID)
	pi := l.view.PrefixIterator(prefix)
	return &iditer3{pi}
}

type iditer3 struct {
	KVIterator
}

func (l *iditer3) TargetID() ID {
	k := l.Key()
	_, _, id3 := split3(k)
	return id3
}

func (l *iditer3) SourceID() ID {
	k := l.Key()
	_, id2, _ := split3(k)
	return id2
}

func (l *iditer3) RelID() ID {
	k := l.Key()
	id1, _, _ := split3(k)
	return id1
}
