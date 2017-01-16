package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"./data"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"golang.org/x/crypto/acme/autocert"
)

var (
	debug = flag.Bool("debug", false, "If set, serve content on HTTP 8080. Otherwise, serve redirects on HTTP 80 and content on HTTPS 443.")
)

const (
	host    = "bran.land"
	email   = "brandon.pitman@gmail.com"
	certDir = "/var/lib/www/certs"
)

type loggingHandler struct {
	h       http.Handler
	logName string
}

func (lh loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s requested %s", lh.logName, r.RemoteAddr, r.RequestURI)
	lh.h.ServeHTTP(w, r)
}

func NewLoggingHandler(logName string, h http.Handler) http.Handler {
	return loggingHandler{
		h:       h,
		logName: logName,
	}
}

type secureHeaderHandler struct {
	h http.Handler
}

func (shh secureHeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src https://fonts.gstatic.com")
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("X-Content-Type-Options", "nosniff")

	shh.h.ServeHTTP(w, r)
}

func NewSecureHeaderHandler(h http.Handler) http.Handler {
	return secureHeaderHandler{
		h: h,
	}
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

func main() {
	flag.Parse()

	// Set up serving mux.
	mux := http.NewServeMux()
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
	handler := NewSecureHeaderHandler(mux)

	// Start serving.
	if *debug {
		server := &http.Server{
			Addr:    "127.0.0.1:8080",
			Handler: NewLoggingHandler("debug", handler),
		}
		log.Printf("Serving debug")
		log.Fatalf("ListenAndServe: %v", server.ListenAndServe())
	}

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
		Handler:      NewLoggingHandler("https", handler),
	}
	log.Printf("Serving")
	go serveHTTPRedirects()
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
