package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	host    = "bran.land"
	email   = "brandon.pitman@gmail.com"
	certDir = "/home/www/certs"
)

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
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			GetCertificate: m.GetCertificate,
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      NewLoggingHandler("https", h),
	}
	log.Printf("Serving")
	go serveHTTPRedirects()
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}

func serveHTTPRedirects() {
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: NewLoggingHandler("http ", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Connection", "close")
			url := *r.URL
			url.Scheme = "https"
			url.Host = host
			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		})),
	}
	log.Fatalf("ListenAndServe: %v", server.ListenAndServe())
}
