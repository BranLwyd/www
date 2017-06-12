// +build debug

package main

import (
	"log"
	"net/http"
)

func serve(h http.Handler) {
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: NewLoggingHandler("debug", h),
	}
	log.Printf("Serving debug")
	log.Fatalf("ListenAndServe: %v", server.ListenAndServe())
}
