package oplog

type OpLog interface {
	Log(interface{}) error
	Scan(count, offset int) ([]interface{}, error)
}

type MemoryOpLog struct {
	log []interface{}
}

func NewMemLog() OpLog {
	return &MemoryOpLog{}
}

func (l *MemoryOpLog) Log(i interface{}) error {
	l.log = append(l.log, i)
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Scan starts from the end
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
