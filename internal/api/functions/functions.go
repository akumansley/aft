package functions

import (
	"context"
	"errors"
	"fmt"

	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

var (
	WrongArity = errors.New("API function called with wrong number of arguments")
	WrongType  = errors.New("API function called with wrong type of argument")
	NoRWTx     = errors.New("API function called without transaction")
)

func unpack(args []interface{}) (ifaceName string, body map[string]interface{}, err error) {
	if len(args) != 2 {
		err = fmt.Errorf("%w: expected %v got %v", WrongArity, 2, args)
		return
	}
	v := args[0]
	ifaceName, ok := v.(string)
	if !ok {
		err = fmt.Errorf("%w: expected %T got %T", WrongType, ifaceName, v)
		return
	}
	v = args[1]
	body, ok = v.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("%w: expected %T got %T", WrongType, body, v)
		return
	}
	return
}

func FindOne(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: tx}
	op, err := p.ParseFindOne(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(tx)
	return
}

var FindOneFunc = db.MakeNativeFunction(
	db.MakeID("43334fbf-015c-41fe-b08d-07edd1dededd"),
	"findOne",
	2,
	db.RPC,
	FindOne,
)

func FindMany(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: tx}
	op, err := p.ParseFindMany(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(tx)
	return
}

var FindManyFunc = db.MakeNativeFunction(
	db.MakeID("31ec796b-7ad5-450c-8eb5-2011672a3f1f"),
	"findMany",
	2,
	db.RPC,
	FindMany,
)

func Count(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: tx}
	op, err := p.ParseCount(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(tx)
	return
}

var CountFunc = db.MakeNativeFunction(
	db.MakeID("afd046c8-8522-4eec-8f08-db561b0adb80"),
	"count",
	2,
	db.RPC,
	Count,
)

func Delete(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseDelete(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var DeleteFunc = db.MakeNativeFunction(
	db.MakeID("55635658-0ed0-43c1-99aa-f1193fc33b6f"),
	"delete",
	2,
	db.RPC,
	Delete,
)

func DeleteMany(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseDeleteMany(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var DeleteManyFunc = db.MakeNativeFunction(
	db.MakeID("8050a482-824d-40ab-9fad-0bbf237ec2c5"),
	"deleteMany",
	2,
	db.RPC,
	DeleteMany,
)

func Update(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}
	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseUpdate(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var UpdateFunc = db.MakeNativeFunction(
	db.MakeID("c4ddc33f-d7e9-4f3a-b95c-4851e0f9c846"),
	"update",
	2,
	db.RPC,
	Update,
)

func UpdateMany(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}

	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseUpdateMany(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var UpdateManyFunc = db.MakeNativeFunction(
	db.MakeID("a6609fd0-8513-42a8-bd8a-6af9ead5554b"),
	"updateMany",
	2,
	db.RPC,
	UpdateMany,
)

func Create(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}

	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseCreate(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var CreateFunc = db.MakeNativeFunction(
	db.MakeID("e475fefd-240c-457c-9b98-4c466e4f674c"),
	"create",
	2,
	db.RPC,
	Create,
)

func Upsert(ctx context.Context, args []interface{}) (result interface{}, err error) {
	modelName, body, err := unpack(args)
	if err != nil {
		return
	}

	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		err = NoRWTx
		return
	}

	p := parsers.Parser{Tx: rwtx}
	op, err := p.ParseUpsert(modelName, body)
	if err != nil {
		return
	}

	eventBus, ok := bus.FromContext(ctx)
	if ok {
		eventBus.Publish(lib.ParseRequest{Request: op})
	}

	result, err = op.Apply(rwtx)
	return
}

var UpsertFunc = db.MakeNativeFunction(
	db.MakeID("6d806ec9-18a8-4fc7-a0e3-2d87835fa337"),
	"upsert",
	2,
	db.RPC,
	Upsert,
)
