package db

import "context"

var FunctionInterface = MakeInterface(
	MakeID("6f55b11e-be7f-4f34-a6ac-1e42d1cd943e"),
	"function",
	[]AttributeL{
		fName,
		fArity,
		fType,
	}, []InterfaceRelationshipL{},
)

func init() {
	FunctionInterface.Relationships_ = []InterfaceRelationshipL{
		FunctionModule,
	}
}

var fName = MakeConcreteAttribute(
	MakeID("048d6151-d80f-44ab-9c77-9ebe70af5b74"),
	"name",
	String,
)

var fArity = MakeConcreteAttribute(
	MakeID("06548160-ce33-427f-9f16-d480253a5c14"),
	"arity",
	Int,
)

var fType = MakeConcreteAttribute(
	MakeID("5ad42d76-cd1f-4fed-a3dc-00d2a1def64f"),
	"funcType",
	FuncType,
)

var FunctionModule = MakeInterfaceRelationship(
	MakeID("c635d40f-a11c-44b1-b96b-3ad714041b89"),
	"module",
	false,
	ModuleModel,
)

type key int

const (
	requestKey key = iota
	txKey
	rwtxKey
)

// defined here because functions are the main client of tx-in-ctx

func WithTx(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func WithRWTx(ctx context.Context, rwtx RWTx) context.Context {
	return context.WithValue(ctx, rwtxKey, rwtx)
}

func RWTxFromContext(ctx context.Context) (rwtx RWTx, ok bool) {
	rwtx, ok = ctx.Value(rwtxKey).(RWTx)
	return
}

func TxFromContext(ctx context.Context) (tx Tx, ok bool) {
	tx, ok = ctx.Value(txKey).(Tx)
	if !ok {
		return RWTxFromContext(ctx)
	}
	return
}
