package api

import (
	"encoding/json"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

func MakeRecord(tx db.Tx, modelName string, jsonValue string) db.Record {
	m, _ := tx.Schema().GetModel(modelName)
	st, err := tx.MakeRecord(m.ID())
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(jsonValue), &st)
	st.Set("id", uuid.New())
	return st
}

type Void struct{}
type Set map[string]Void
