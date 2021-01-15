package db

type UnionNode struct {
	nodes []Node
}

func (u *UnionNode) String() string {
	return "UnionNode{}"
}

func (u *UnionNode) Children() []Node {
	return u.nodes
}

func (u *UnionNode) ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error) {
	return &unionIterator{tx: tx, qr: qr, nodes: u.nodes, returned: map[ID]*QueryResult{}}, nil
}

type unionIterator struct {
	tx       *holdTx
	qr       *QueryResult
	nodes    []Node
	current  qrIterator
	returned map[ID]*QueryResult
	value    *QueryResult
	pos      int
	err      error
}

func (i *unionIterator) setIterator() (err error) {
	i.current, err = i.nodes[i.pos].ResultIter(i.tx, i.qr)
	return
}

func (i *unionIterator) Next() bool {
	for {
		if i.current == nil {
			err := i.setIterator()
			if err != nil {
				i.err = err
				return false
			}
		}
		ok := i.current.Next()
		if ok {
			qr := i.current.Value()
			// if we've already returned this, skip it
			if _, ok := i.returned[qr.Record.ID()]; ok {
				continue
			}
			i.returned[qr.Record.ID()] = qr
			i.value = qr
			return true
		}

		if i.current.Err() == Done && i.pos+1 < len(i.nodes) {
			i.pos++
			err := i.setIterator()
			if err != nil {
				i.err = err
				return false
			}
		} else {
			i.err = i.current.Err()
			return false
		}
	}
}

func (i *unionIterator) Value() *QueryResult {
	return i.value
}

func (i *unionIterator) Err() error {
	return i.err
}
