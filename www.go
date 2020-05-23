package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
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
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("X-XSS-Protection", "1; mode=block")
	w.Header().Add("X-Content-Type-Options", "nosniff")

	// Chrome's PDF renderer uses inline CSS, which is broken by strict CSPs (!).
	// TODO: remove the 'unsafe-inline' style-src directive once Chrome fixes its PDF renderer
	w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com")

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

	tagOnce sync.Once
	tag     string
}

func (sh *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", sh.contentType)
	w.Header().Set("ETag", sh.etag())
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(sh.content))
}

func (sh *staticHandler) etag() string {
	sh.tagOnce.Do(func() {
		h := sha256.Sum256(sh.content)
		sh.tag = fmt.Sprintf(`"%s"`, base64.RawURLEncoding.EncodeToString(h[:]))
	})
	return sh.tag
}

func NewAssetHandler(assetName, contentType string) (*staticHandler, error) {
	content, ok := assets.Asset[assetName]
	if !ok {
		return nil, fmt.Errorf("no such asset %q", assetName)
	}
	sh := &staticHandler{
		content:     content,
		contentType: contentType,
	}
	go sh.etag() // eagerly compute etag so that it will probably be available by the first request
	return sh, nil
}

func Must(h http.Handler, err error) http.Handler {
	if err != nil {
		panic(err)
	}
	return h
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", NewFilteredHandler("/", Must(NewAssetHandler("assets/index.html", "text/html; charset=utf-8"))))
	mux.Handle("/flowers", flowerHandler{})
	mux.Handle("/style.css", Must(NewAssetHandler("assets/style.css", "text/css; charset=utf-8")))
	mux.Handle("/favicon.ico", Must(NewAssetHandler("assets/favicon.ico", "image/x-icon")))
	mux.Handle("/resume.pdf", Must(NewAssetHandler("assets/resume.pdf", "application/pdf")))

	// Flower pics.
	for _, species := range []string{"R", "T", "P", "C", "L", "H", "W", "M"} {
		for _, color := range []string{"W", "P", "R", "O", "Y", "G", "B", "U", "K"} {
			assetName := fmt.Sprintf("assets/img/%s%s.png", species, color)
			if _, ok := assets.Asset[assetName]; !ok {
				// Not all species/color combinations are represented.
				continue
			}
			mux.Handle(fmt.Sprintf("/img/%s%s.png", species, color), Must(NewAssetHandler(assetName, "image/png")))
		}
	}

	serve(NewSecureHeaderHandler(mux))
}
