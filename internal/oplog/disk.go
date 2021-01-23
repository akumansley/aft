package oplog

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/awans/logio/logio"
)

type LogStore interface {
	Iterator() ByteIterator
	Log([]byte) error
	Scan(count, offset int) ([][]byte, error)
	Close()
}

type OpLog interface {
	Iterator() Iterator
	Log(interface{}) error
	Scan(count, offset int) ([]interface{}, error)
	Close()
}

type Iterator interface {
	Next() bool
	Value() interface{}
	Err() error
}

type ByteIterator interface {
	Next() bool
	Value() []byte
	Err() error
}

type DiskOpLogIterator struct {
	log   *DiskOpLog
	off   int64
	value []byte
	err   error
	done  bool
}

func (i *DiskOpLogIterator) Value() []byte {
	if i.done {
		panic("Called value after done")
	}
	return i.value
}
func (i *DiskOpLogIterator) Err() error {
	return i.err
}

func (i *DiskOpLogIterator) Next() bool {
	i.log.Lock()
	defer i.log.Unlock()
	if i.log.closed {
		panic("closed")
	}
	if i.off > i.log.tail {
		i.done = true
		return false
	}

	i.log.f.Seek(i.off, io.SeekStart)
	reader := logio.NewReader(i.log.f, i.off)
	bts, err := reader.Read()
	rlen := len(bts)

	if err != nil {
		if err == io.EOF {
			i.done = true
			return false
		}
		i.err = err
		return false
	}

	// 15 is the length in bytes of the record header
	i.off += int64(rlen) + 15
	i.value = bts
	return true
}

func OpenDiskLog(filename string) (LogStore, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	off := info.Size()
	tail, err := logio.Rewind(f, off)
	if err != nil {
		// eof is fine
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	return &DiskOpLog{f: f, tail: tail}, nil
}

type DiskOpLog struct {
	sync.Mutex
	f      *os.File
	tail   int64
	closed bool
}

func (l *DiskOpLog) Close() {
	l.Lock()
	l.f.Close()
	l.closed = true
	l.Unlock()
}

func (l *DiskOpLog) Log(entry []byte) error {
	l.Lock()
	defer l.Unlock()
	if l.closed {
		panic("closed")
	}

	lw := logio.NewWriter(l.f, l.tail)

	err := lw.Append(entry)
	if err != nil {
		return err
	}
	l.tail = lw.Tell()

	err = l.f.Sync()
	return err
}

// Scan starts from the end and goes backwards
func (l *DiskOpLog) Scan(count, offset int) ([][]byte, error) {
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
	var entries [][]byte
	for {
		off, err = logio.Rewind(l.f, off)
		if err == io.EOF {
			break
		}
		l.f.Seek(off, io.SeekStart)
		entry, err := logio.NewReader(l.f, off).Read()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		count--
		if count == 0 {
			break
		}
	}
	return entries, nil
}

func (l *DiskOpLog) Iterator() ByteIterator {
	return &DiskOpLogIterator{log: l, off: 0}
}
