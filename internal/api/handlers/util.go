package handlers

import (
	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}

type SummaryResponse struct {
	Count int `json:"count"`
}

func unpackArgs(r *http.Request) (string, map[string]interface{}, error) {
	var body map[string]interface{}
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	buf, _ := ioutil.ReadAll(r.Body)
	err := jsoniter.Unmarshal(buf, &body)
	if err != nil {
		return "", body, err
	}
	return modelName, body, nil
}

func response(w http.ResponseWriter, result interface{}) {
	bytes, _ := jsoniter.Marshal(&result)
	_, _ = w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}
