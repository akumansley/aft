package server

import (
	"awans.org/aft/internal/server/db"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run() {
	db.SetupTestData()
	r := NewRouter()
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
