package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"./data"

	"golang.org/x/crypto/acme/autocert"
)

const (
	host    = "bran.land"
	email   = "brandon.pitman@gmail.com"
	certDir = "/var/lib/www/certs"
)

var (
	debug = flag.Bool("debug", false, "If set, serve content on HTTP 8080. Otherwise, serve redirects on HTTP 80 and content on HTTPS 443.")

	dataIndex   = data.MustAsset("index.html")
	dataStyle   = data.MustAsset("style.css")
	dataFavicon = data.MustAsset("favicon.ico")
)

type loggingHandler struct {
	h       http.Handler
	logName string
}

func (lh loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip port from remote address, as the client port is not useful information.
	ra := r.RemoteAddr
	idx := strings.LastIndex(ra, ":")
	if idx != -1 {
		ra = ra[:idx]
	}
	log.Printf("[%s] %s requested %s", lh.logName, ra, r.RequestURI)
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

// filteredHandler filters a handler to only serve one path; anything else is given a 404.
type filteredHandler struct {
	allowedPath string
	h           http.Handler
}

func (fh filteredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != fh.allowedPath {
		http.NotFound(w, r)
	} else {
		fh.h.ServeHTTP(w, r)
	}
}

func NewFilteredHandler(allowedPath string, h http.Handler) http.Handler {
	return &filteredHandler{
		allowedPath: allowedPath,
		h:           h,
	}
}

// staticHandler serves static content from memory.
type staticHandler struct {
	content     []byte
	contentType string
}

func (sh staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", sh.contentType)
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(sh.content))
}

func NewStaticHandler(content []byte, contentType string) *staticHandler {
	return &staticHandler{
		content:     content,
		contentType: contentType,
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
	mux.Handle("/", NewFilteredHandler("/", NewStaticHandler(dataIndex, "text/html; charset=utf-8")))
	mux.Handle("/style.css", NewStaticHandler(dataStyle, "text/css; charset=utf-8"))
	mux.Handle("/favicon.ico", NewStaticHandler(dataFavicon, "image/x-icon"))
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
