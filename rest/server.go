package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const port = ":3000"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/salesReport", SalesReportHandler)
	log.Println("Starting server on address", port)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
