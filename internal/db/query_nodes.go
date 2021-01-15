package db

import (
	"fmt"

	"github.com/disiqueira/gotree"
)

type Node interface {
	fmt.Stringer
	Children() []Node
	ResultIter(tx *holdTx, qr *QueryResult) (qrIterator, error)
}

type qrIterator interface {
	Next() bool
	Value() *QueryResult
	Err() error
}

func PrintTree(n Node) {
	t := gotree.New(n.String())
	for _, child := range n.Children() {
		PrintTreeRec(t, child)
	}
	fmt.Println(t.Print())
}

func PrintTreeRec(t gotree.Tree, n Node) {
	if n == nil {
		t = t.Add("Nil")
	} else {
		t = t.Add(n.String())
		for _, child := range n.Children() {
			PrintTreeRec(t, child)
		}
	}
}
