package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const (
	host    = "bran.land"
	email   = "brandon.pitman@gmail.com"
	certDir = "/home/www/certs"
)

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1") // enable TLS 1.3; remove once enabled by default
}

func serve(h http.Handler) {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache(certDir),
		Email:      email,
	}
	server := &http.Server{
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			MinVersion:             tls.VersionTLS13,
			SessionTicketsDisabled: true,
			GetCertificate:         m.GetCertificate,
			NextProtos:             []string{"h2", acme.ALPNProto},
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      NewLoggingHandler("https", h),
	}

	log.Printf("Serving")
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
