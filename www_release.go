package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
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
				tls.X25519,
				tls.CurveP256,
			},
			MinVersion:             tls.VersionTLS12,
			SessionTicketsDisabled: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			GetCertificate: m.GetCertificate,
			NextProtos:     []string{"h2", acme.ALPNProto},
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      NewLoggingHandler("https", h),
	}

	log.Printf("Serving")
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
