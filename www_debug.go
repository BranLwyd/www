package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

var (
	hostname = flag.String("hostname", "", "The hostname to serve on. Defaults to os.Hostname().")
)

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1") // enable TLS 1.3; remove once enabled by default
}

func serve(h http.Handler) {
	// Figure out hostname if it is not set.
	if *hostname == "" {
		hn, err := os.Hostname()
		if err != nil {
			log.Fatalf("Could not determine hostname (and not overridden by --hostname): %v", err)
		}
		*hostname = hn
	}

	// Generate a self-signed certificate with the appropriate hostname.
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Could not generate key: %v", err)
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Harpocrates"},
		},
		DNSNames:              []string{*hostname},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, priv.Public(), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}

	// Begin serving.
	server := &http.Server{
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{{
				Certificate: [][]byte{certDER},
				PrivateKey:  priv,
				Leaf:        cert,
			}},
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			MinVersion:             tls.VersionTLS13,
			SessionTicketsDisabled: true,
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":8080",
		Handler:      NewLoggingHandler("debug", h),
	}
	log.Printf("Serving debug on %s:8080", *hostname)
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
