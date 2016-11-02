package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"./data"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

const (
	host           = "bran.land"
	certFile       = "/var/lib/www/cert.crt"
	keyFile        = "/var/lib/www/decrypted.key"
	leChallengeDir = "/var/lib/www/.le/.well-known"
)

var (
	mux *http.ServeMux
)

type certCache struct {
	certFile string
	keyFile  string

	certMu sync.RWMutex
	cert   *tls.Certificate
}

func newCertCache(certFile string, keyFile string, refreshInterval time.Duration) (*certCache, error) {
	cc := &certCache{
		certFile: certFile,
		keyFile:  keyFile,
	}

	if err := cc.set(); err != nil {
		return nil, err
	}

	go func() {
		for range time.Tick(refreshInterval) {
			log.Print("Reloading certificate")
			if err := cc.set(); err != nil {
				log.Printf("Could not reload certificate: %v", err)
			}
		}
	}()

	return cc, nil
}

func (cc *certCache) set() error {
	cert, err := tls.LoadX509KeyPair(cc.certFile, cc.keyFile)
	if err != nil {
		return err
	}

	cc.certMu.Lock()
	defer cc.certMu.Unlock()
	cc.cert = &cert
	return nil
}

func (cc *certCache) Get() *tls.Certificate {
	cc.certMu.RLock()
	defer cc.certMu.RUnlock()
	return cc.cert
}

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
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectHandler)
	mux.Handle("/.well-known/", http.StripPrefix("/.well-known/", http.FileServer(http.Dir(leChallengeDir))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: NewLoggingHandler("http ", mux),
	}
	log.Fatalf("ListenAndServe: %v", server.ListenAndServe())
}

func main() {
	go httpHandler()

	// Read cert.
	log.Printf("Loading certificate")
	certCache, err := newCertCache("/var/lib/www/cert.crt", "/var/lib/www/decrypted.key", time.Hour)
	if err != nil {
		log.Fatalf("Could not load certificate: %v", err)
	}

	// Set up serving mux.

	// Start serving.
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
		MinVersion: tls.VersionTLS10,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		},
		Certificates:   []tls.Certificate{*certCache.Get()}, // This will never be used, but is required.
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) { return certCache.Get(), nil },
	}
	server := &http.Server{
		Addr:      ":10443",
		Handler:   NewLoggingHandler("https", http.HandlerFunc(contentHandler)),
		TLSConfig: config,
	}
	log.Printf("Serving")
	log.Fatalf("ListenAndServeTLS: %v", server.ListenAndServeTLS("", ""))
}
