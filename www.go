package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"./data"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"golang.org/x/crypto/acme/autocert"
)

const (
	host    = "bran.land"
	email   = "brandon.pitman@gmail.com"
	certDir = "/var/lib/www/certs"
)

var (
	mux *http.ServeMux
)

type loggingHandler struct {
	h       http.Handler
	logName string
}

func (lh *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s requested %s", lh.logName, r.RemoteAddr, r.RequestURI)
	lh.h.ServeHTTP(w, r)
}

func NewLoggingHandler(logName string, h http.Handler) *loggingHandler {
	return &loggingHandler{
		h:       h,
		logName: logName,
	}
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src https://fonts.gstatic.com")
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("X-Content-Type-Options", "nosniff")

	mux.ServeHTTP(w, r)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	url := *r.URL
	url.Scheme = "https"
	url.Host = host
	http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
}

func httpHandler() {
	server := &http.Server{
		Addr:    ":http",
		Handler: NewLoggingHandler("http ", http.HandlerFunc(redirectHandler)),
	}
	log.Fatalf("ListenAndServe: %v", server.ListenAndServe())
}

func main() {
	go httpHandler()

	// Set up certificate handling.
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),
		Cache:      autocert.DirCache(certDir),
		Email:      email,
	}

	// Set up serving mux.
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: data.Asset, AssetDir: data.AssetDir, AssetInfo: data.AssetInfo}))
	mux.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		idx := strings.Index(r.RemoteAddr, ":")
		if idx == -1 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ip := r.RemoteAddr[:idx]

		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, ip)
	})

	// Start serving.
	config := &tls.Config{
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
	}
	server := &http.Server{
		Addr:      ":https",
		Handler:   NewLoggingHandler("https", http.HandlerFunc(contentHandler)),
		TLSConfig: config,
	}
	log.Printf("Serving")
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
