package oplog

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/awans/logio/logio"
	"io"
	"os"
	"sync"
)

type OpLog interface {
	Iterator() Iterator
	Log(interface{}) error
	Scan(count, offset int) ([]interface{}, error)
	Close()
}

type Iterator interface {
	Next() (interface{}, bool)
}

type GobLogIterator struct {
	log *GobLog
	off int64
}

func (i *GobLogIterator) Next() (interface{}, bool) {
	i.log.Lock()
	defer i.log.Unlock()
	if i.log.closed {
		panic("closed")
	}

	i.log.f.Seek(i.off, io.SeekStart)
	bts, err := logio.NewReader(i.log.f, i.off).Read()
	rlen := len(bts)
	if err != nil {
		return nil, false
	}
	buf := bytes.NewBuffer(bts)
	var entry interface{}
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&entry)
	if err != nil {
		panic(err)
	}
	i.off += int64(rlen)
	return entry, true
}

func initGob() {
	gob.Register(CreateOp{})
	gob.Register(ConnectOp{})
	gob.Register(TxEntry{})
}

func OpenGobLog(filename string) (OpLog, error) {
	initGob()
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("createerror\n")
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		fmt.Printf("staterr\n")
		return nil, err
	}
	off := info.Size()
	tail, err := logio.Rewind(f, off)
	if err != nil {
		// eof is fine
		if !errors.Is(err, io.EOF) {
			fmt.Printf("rewinderr\n")
			return nil, err
		}
	}
	return &GobLog{f: f, tail: tail}, nil
}

type GobLog struct {
	sync.Mutex
	f      *os.File
	tail   int64
	closed bool
}

func (l *GobLog) Close() {
	l.Lock()
	l.f.Close()
	l.closed = true
	l.Unlock()
}

func (l *GobLog) Log(i interface{}) error {
	l.Lock()
	defer l.Unlock()
	if l.closed {
		panic("closed")
	}

	fmt.Printf("log:%+v\n", i)

	lw := logio.NewWriter(l.f, l.tail)
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)
	err := e.Encode(&i)
	if err != nil {
		fmt.Printf("goberr: %v\n", err)
		return err
	}
	bts := buf.Bytes()
	err = lw.Append(bts)
	if err != nil {
		fmt.Printf("logio.append err: %v\n", err)
		return err
	}
	l.tail = lw.Tell()

	// TODO make this happen less often?
	l.f.Sync()
	return err
}

// Scan starts from the end and goes backwards
func (l *GobLog) Scan(count, offset int) ([]interface{}, error) {
	l.Lock()
	defer l.Unlock()
	if l.closed {
		panic("closed")
	}

	info, err := l.f.Stat()
	if err != nil {
		return nil, err
	}
	off := info.Size()
	var entries []interface{}
	for {
		off, err = logio.Rewind(l.f, off)
		if err == io.EOF {
			break
		}
		l.f.Seek(off, io.SeekStart)
		bts, err := logio.NewReader(l.f, off).Read()
		if err != nil {
			return nil, err
		}
		var entry interface{}
		buf := bytes.NewBuffer(bts)
		dec := gob.NewDecoder(buf)
		dec.Decode(&entry)
		entries = append(entries, entry)
		count--
		if count == 0 {
			break
		}
	}
	return entries, nil
}

func (l *GobLog) Iterator() Iterator {
	return &GobLogIterator{log: l, off: 0}
}