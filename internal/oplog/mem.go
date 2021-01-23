package oplog

type ApiOpEntry struct {
	OpType int
	body   interface{}
}

type MemoryOpLog struct {
	log [][]byte
}

type MemoryOpLogIterator struct {
	log   *MemoryOpLog
	ix    int
	value []byte
	err   error
}

func (i *MemoryOpLogIterator) Value() []byte {
	return i.value
}
func (i *MemoryOpLogIterator) Err() error {
	return i.err
}

func (i *MemoryOpLogIterator) Next() bool {
	if i.ix < len(i.log.log) {
		i.ix++
		i.value = i.log.log[i.ix-1]
		return true
	}
	return false
}

func NewMemLog() LogStore {
	return &MemoryOpLog{}
}

func (l *MemoryOpLog) Log(i []byte) error {
	l.log = append(l.log, i)
	return nil
}

func (l *MemoryOpLog) Close() {
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Scan starts from the end and goes backwards
func (l *MemoryOpLog) Scan(count, offset int) ([][]byte, error) {
	startIx := len(l.log) - 1 - offset

	// return what you can even if it doesn't fill count
	toIx := max(startIx-count, 0)
	if startIx < 0 {
		return [][]byte{}, nil
	}

	var resp [][]byte
	for i := startIx; i >= toIx; i-- {
		resp = append(resp, l.log[i])
	}
	return resp, nil
}

func (l *MemoryOpLog) Iterator() ByteIterator {
	return &MemoryOpLogIterator{log: l, ix: 0}
}
