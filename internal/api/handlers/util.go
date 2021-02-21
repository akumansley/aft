package handlers

import (
	"io/ioutil"
	"net/http"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/db"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}

type SummaryResponse struct {
	Count int `json:"count"`
}

func unpackArgs(r *http.Request) (methodName, modelName string, body map[string]interface{}, err error) {
	vars := mux.Vars(r)
	methodName = vars["methodName"]
	modelName = vars["modelName"]
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(buf, &body)
	return
}

func response(w http.ResponseWriter, result interface{}) {
	bytes, _ := jsoniter.Marshal(&result)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}
func AddFunctionLiterals(testDB db.DB) {
	funcs := []db.NativeFunctionL{
		functions.FindOneFunc,
		functions.FindManyFunc,
		functions.CountFunc,
		functions.DeleteFunc,
		functions.DeleteManyFunc,
		functions.UpdateFunc,
		functions.UpdateManyFunc,
		functions.CreateFunc,
		functions.UpsertFunc,
	}
	rwtx := testDB.NewRWTx()
	for _, f := range funcs {
		testDB.AddLiteral(rwtx, f)
		testDB.RegisterNativeFunction(f)
	}
	rwtx.Commit()
}
