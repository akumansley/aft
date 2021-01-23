package audit

import (
	"errors"
	"fmt"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/oplog"
)

type LogScanHandler struct {
	log oplog.OpLog
	bus *bus.EventBus
}

var ErrArgs = errors.New("argument-error")

func makeScanFunction(logs map[string]oplog.OpLog) db.FunctionL {
	scanFunc := func(args []interface{}) (result interface{}, err error) {
		input := args[1]

		rpcData := input.(map[string]interface{})

		logVal, ok := rpcData["log"]
		if !ok {
			return nil, fmt.Errorf("%w: log is required", ErrArgs)
		}

		logName, ok := logVal.(string)
		if !ok {
			return nil, fmt.Errorf("%w: log expected string got %T", ErrArgs, logVal)
		}

		countVal, ok := rpcData["count"]
		if !ok {
			return nil, fmt.Errorf("%w: count is required", ErrArgs)
		}
		count, ok := countVal.(float64)
		if !ok {
			return nil, fmt.Errorf("%w: expected int got %T", ErrArgs, countVal)
		}

		offsetVal, ok := rpcData["offset"]
		if !ok {
			return nil, fmt.Errorf("%w: offset is required", ErrArgs)
		}
		offset, ok := offsetVal.(float64)
		if !ok {
			return nil, fmt.Errorf("%w: expected int got %T", ErrArgs, offsetVal)
		}

		log, ok := logs[logName]
		if !ok {
			return nil, errors.New("Invalid logName")
		}

		v, err := log.Scan(int(count), int(offset))
		return v, err

	}

	var scanRPC = db.MakeNativeFunction(
		db.MakeID("8518fcea-8826-409d-8274-3e4373a4d971"),
		"scan",
		2,
		scanFunc,
	)

	return scanRPC
}
