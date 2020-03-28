package server

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/server/services"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var objects = services.ListObjectsResponse{
	Data: []*data.Object{
		&data.Object{
			Id:   "Cekw67uyMpBGZLRP2HFVbe",
			Name: "Test",
			Fields: []*data.Field{
				&data.Field{
					Name: "f1",
					Type: data.FieldType_TEXT,
				},
			},
		},
		&data.Object{
			Id:   "6R7VqaQHbzC1xwA5UueGe6",
			Name: "Cool",
			Fields: []*data.Field{
				&data.Field{
					Name: "f5",
					Type: data.FieldType_INT,
				},
			},
		},
	},
}

func InfoObject(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var params services.InfoObjectsRequest
	buf, _ := ioutil.ReadAll(req.Body)
	_ = protojson.Unmarshal(buf, &params)
	id := params.Id
	var object *data.Object
	for _, obj := range objects.Data {
		if obj.Id == id {
			object = obj
			break
		}
	}
	if object != nil {
		response := services.InfoObjectsResponse{
			Object: object,
		}
		bytes, _ := protojson.Marshal(&response)
		_, _ = w.Write(bytes)
	} else {
		http.NotFound(w, req)
	}
}

func ListObjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	bytes, err := protojson.Marshal(&objects)
	if err != nil {
		http.NotFound(w, req)
	}
	_, err = w.Write(bytes)
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
