package server

import (
	"awans.org/aft/internal/data"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var objects = []data.Object{
	data.Object{
		Id:   "Cekw67uyMpBGZLRP2HFVbe",
		Name: "Test",
		Fields: []*data.Field{
			&data.Field{
				Name: "f1",
				Type: data.FieldType_TEXT,
			},
		},
	},
	data.Object{
		Id:   "6R7VqaQHbzC1xwA5UueGe6",
		Name: "Cool",
		Fields: []*data.Field{
			&data.Field{
				Name: "f5",
				Type: data.FieldType_INT,
			},
		},
	},
}

type InfoObjectParams struct {
	Id string `json:"id"`
}

func InfoObject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var params InfoObjectParams
	json.NewDecoder(req.Body).Decode(&params)
	id := params.Id
	var object *data.Object
	for _, obj := range objects {
		if obj.Id == id {
			object = &obj
			break
		}
	}
	if object != nil {
		json.NewEncoder(w).Encode(object)
	} else {
		http.NotFound(w, req)
	}
}

func ListObjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(objects)
}

func Run() {
	r := mux.NewRouter()
	s := r.Methods("POST").Subrouter()
	s.HandleFunc("/api/objects.info", InfoObject)
	s.HandleFunc("/api/objects.list", ListObjects)
	port := ":8080"
	fmt.Println("Serving on port", port)
	srv := &http.Server{
		Handler: r,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
}
