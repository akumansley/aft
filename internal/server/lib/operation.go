package lib

type TxType uint

const (
	Tx TxType = iota
	RWTx
	None
)

type Operation struct {
	Name    string
	Service string
	Method  string
	Server  Server
	Tx      TxType
}
