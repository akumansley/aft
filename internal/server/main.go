package server
import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"awans.org/aft/internal/data"
)

func ListObjects(w http.ResponseWriter, req *http.Request) {
	objects := []data.Object{
		data.Object{
			Name: "Test",
			Fields: []data.Field{
				data.Field{
					Name: "f1",
					Type: data.Text,
				},
			},
		},
		data.Object{
			Name: "Cool",
			Fields: []data.Field{
				data.Field{
					Name: "f5",
					Type: data.Int,
				},
			},
		},
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(objects)
}

func Run() {
	http.HandleFunc("/objects", ListObjects)
	port := ":8080"
	fmt.Println("Serving on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
