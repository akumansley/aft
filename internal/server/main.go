package server

import (
	"awans.org/aft/internal/data"
	"awans.org/aft/internal/server/services"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/steveyen/gtreap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"math/rand"
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
	object = &data.Object{}
	objBuf := objectTable.Get(id)
	proto.Unmarshal(objBuf, object)
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

type item struct {
	id    string
	bytes []byte // this.. doesn't make sense
	// we should just store a pointer to an interface{}
	// or switch to using flatbuffers
}

type Table struct {
	t *gtreap.Treap
}

func stringIdCompare(a, b interface{}) int {
	return bytes.Compare([]byte(a.(item).id), []byte(b.(item).id))
}

func (t *Table) Init() {
	t.t = gtreap.NewTreap(stringIdCompare)
}

func (t *Table) Upsert(id string, bytes []byte) {
	t.t = t.t.Upsert(item{id: id, bytes: bytes}, rand.Int()) // rand approximates balanced
}

func (t *Table) Get(id string) []byte {
	it := t.t.Get(item{id: id})
	return it.(item).bytes
}

var objectTable Table

func setupTestData() {
	objectTable = Table{}
	objectTable.Init()
	for _, obj := range objects.Data {
		bytes, _ := proto.Marshal(obj)
		objectTable.Upsert(obj.Id, bytes)
	}
}

func Run() {
	setupTestData()

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
