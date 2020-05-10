package oplog

type OpLog interface {
	Iterator() Iterator
	Log(interface{}) error
	Scan(count, offset int) ([]interface{}, error)
	NextId() uint
}

type Iterator interface {
	Next() (interface{}, bool)
}

type ApiOpEntry struct {
	OpId   int
	OpType int
	body   interface{}
}

type MemoryOpLog struct {
	log []interface{}
}

type MemoryOpLogIterator struct {
	log *MemoryOpLog
	ix  int
}

func (i *MemoryOpLogIterator) Next() (interface{}, bool) {
	if i.ix < len(i.log.log) {
		i.ix++
		return i.log.log[i.ix-1], true
	}
	return nil, false
}

func NewMemLog() OpLog {
	return &MemoryOpLog{}
}

func (l *MemoryOpLog) Log(i interface{}) error {
	l.log = append(l.log, i)
	return nil
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Scan starts from the end and goes backwards
func (l *MemoryOpLog) Scan(count, offset int) ([]interface{}, error) {
	startIx := len(l.log) - 1 - offset

	// return what you can even if it doesn't fill count
	toIx := max(startIx-count, 0)
	if startIx < 0 {
		return []interface{}{}, nil
	}

	var resp []interface{}
	for i := startIx; i >= toIx; i-- {
		resp = append(resp, l.log[i])
	}
	return resp, nil
}

func (l *MemoryOpLog) Iterator() Iterator {
	return &MemoryOpLogIterator{log: l, ix: 0}
}

func (l *MemoryOpLog) NextId() uint {
	return uint(len(l.log))
}
