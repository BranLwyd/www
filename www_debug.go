package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	hostname = flag.String("hostname", "", "The hostname to serve on. Defaults to os.Hostname().")
)

func serve(h http.Handler) {
	log.Printf("Serving debug on :8080")
	log.Fatalf("ListenAndServe: %v", http.ListenAndServe(":8080", NewLoggingHandler("debug", h)))

}
