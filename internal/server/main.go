package server

import (
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/services"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Run() {
	db.SetupTestData()

	r := mux.NewRouter()
	s := r.Methods("POST").Subrouter()
	s.HandleFunc("/api/objects.info", services.InfoObjects)
	s.HandleFunc("/api/objects.list", services.ListObjects)
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
