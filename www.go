package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BranLwyd/www/assets"
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
	w.Header().Add("Content-Security-Policy", "default-src 'self'")
	w.Header().Add("X-Frame-Options", "DENY")
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
	modTime     time.Time
}

func (sh staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", sh.contentType)
	http.ServeContent(w, r, "", sh.modTime, bytes.NewReader(sh.content))
}

func NewAssetHandler(assetName, contentType string) (*staticHandler, error) {
	content, err := assets.Asset(assetName)
	if err != nil {
		return nil, fmt.Errorf("could not get asset %q: %v", assetName, err)
	}
	info, err := assets.AssetInfo(assetName)
	if err != nil {
		return nil, fmt.Errorf("could not get asset info for %q: %v", assetName, err)
	}
	return &staticHandler{
		content:     content,
		contentType: contentType,
		modTime:     info.ModTime(),
	}, nil
}

func Must(h http.Handler, err error) http.Handler {
	if err != nil {
		panic(err)
	}
	return h
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", NewFilteredHandler("/", Must(NewAssetHandler("index.html", "text/html; charset=utf-8"))))
	mux.Handle("/style.css", Must(NewAssetHandler("style.css", "text/css; charset=utf-8")))
	mux.Handle("/favicon.ico", Must(NewAssetHandler("favicon.ico", "image/x-icon")))
	serve(NewSecureHeaderHandler(mux))
}
