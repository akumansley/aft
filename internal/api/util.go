package api

import (
	"awans.org/aft/internal/db"
	"encoding/json"
)

func MakeRecord(tx db.Tx, modelName string, jsonValue string) db.Record {
	m, _ := tx.Schema().GetModel(modelName)
	st, err := tx.MakeRecord(m.ID())
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}
