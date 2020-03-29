package server

import (
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/services"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run() {
	db.SetupTestData()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/objects.info", services.InfoObjects)
	mux.HandleFunc("/api/objects.list", services.ListObjects)
	port := ":8080"
	fmt.Println("Serving on port", port)
	srv := &http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
}
