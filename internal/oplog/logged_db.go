package oplog

import (
	"awans.org/aft/internal/db"
)

func DBFromLog(appDB db.DB, l OpLog) error {
	iter := l.Iterator()
	for iter.Next() {
		rwtx := appDB.NewRWTx()
		val := iter.Value()
		opList := val.([]db.Operation)
		for _, op := range opList {
			op.Replay(rwtx)
		}
		err := rwtx.Commit()
		if err != nil {
			return err
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	return nil
}

func MakeTransactionLogger(l OpLog) func(db.BeforeCommit) {
	logger := func(event db.BeforeCommit) {
		ops := event.Tx.Operations()

		if len(ops) == 0 {
			return
		}
		err := l.Log(ops)
		if err != nil {
			panic(err)
		}
	}
	return logger
}
