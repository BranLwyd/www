package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	host           = "bran.land"
	certFile       = "/var/lib/www/cert.crt"
	keyFile        = "/var/lib/www/decrypted.key"
	leChallengeDir = "/var/lib/www/.le/.well-known"
)

var (
	favicon string
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

	switch r.URL.Path {
	case "/":
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, page)

	case "/style.css":
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprint(w, style)

	case "/favicon.ico":
		w.Header().Add("Content-Type", "image/png")
		fmt.Fprint(w, favicon)

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
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
	var err error
	faviconBytes, err := base64.StdEncoding.DecodeString(faviconBase64)
	if err != nil {
		log.Fatalf("Could not initialize favicon: %v", err)
	}
	favicon = string(faviconBytes)

	go httpHandler()

	// Read cert.
	log.Printf("Loading certificate")
	certCache, err := newCertCache("/var/lib/www/cert.crt", "/var/lib/www/decrypted.key", time.Hour)
	if err != nil {
		log.Fatalf("Could not load certificate: %v", err)
	}

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

// Website data below here.
const (
	faviconBase64 = `
iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAQAAAAAYLlVAAACmklEQVR4Ae3YU5DEaBiG0TPm2rZt
28bN2rZt27Zt27Zt2+a3TlUPkunUn/V5rttI8vrfP0Or1d3veyF86DabGsefZGInesFbojZfetxE
KtNmEgvY3csit1VV4gA/+M6PQnGbSW4fUVeTqcSEVvem7/0gCnpYhYY1q529I3KbVuXeFzmdp3Jb
iJye065yN4t++8HQKjeByGlJlRtH5PSj+VRsCVHQNiq1jCjsOBU6QhTnSW0q8rkYUE/rVIGVxYB7
w0iSu1XU0RdGkdhHoq7u+qufQFhJQg0+FnU33F/9BI6XTLdPRYlG/aufwMsSafSJKNVCEnlWlOo+
iewqSraMJIYXJXtXmyTuFiVbVBJLi5J9JJG9Rcn2k8hIbhEl+tDgEun0nCjRARJ6XJQooelFiRaQ
TJcPRN2dKaHDRN3dpEE6JQ5On2uS0Dx+EHW1ggZJTVbzLpxTeLg6VXKN1vr9SbwEdhK5LSa5JncJ
Ye/8DSFrgqouWTsAa4vcrpLcVMKMAMYTBS0osVucK4OXRG4f6JLQ4kKLDPYXBa0rmQ7fOlStC4Vj
RW6jSWQjP5parXCuGQoXhETe9bpaCwibGDF3zLjDTJLYWpijj5l7Vp3eEP32oCSG8JWr1eJyXxsP
94ickthW2M1Qaj3r+eJh11zKQwPmko20ZgPQIjyuCXwr+u12pQzmOFv93GUeyMaodc1nZS2A44Uv
jAk2FP10vtWVcowNPZy97m98ZwVAiw6D2UQIH2sHnT4Uffa9krpwrfCdY81kElNoAkO7z7M+zu5+
usJxe2slNThTb4f6zLvCB9ZwgDArgOV8IfroHaU16+1801rL+4YFVxlPxiiek332zhC/l9TImi2g
EzTobQ2XOtd4YGKHusO1/lfjJ/sX7sqBgKuSAAAAAElFTkSuQmCC`

	page = `
<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=0.5">
	<title>Branland</title>
	<link rel='stylesheet' type='text/css' href='https://fonts.googleapis.com/css?family=Droid+Serif'>
	<link rel="stylesheet" type='text/css' href="style.css">
</head>
<body>
	<div class="ribbon"></div>
	<div class="content">
		<div class="header">
			<h1>branland</h1>
		</div>

		<div class="inner-content">
			<h2>Introduction</h2>
			<p>Hi, I'm Brandon. Welcome to Branland! This is my personal website.  I'm a Software Engineer in the Bay Area. My interests include mathematics, computer science, and information security. I'm also into board games, tabletop role-playing games, science fiction, coffee, wine, and whiskey. I'm relatively new to the Bay, so I'm always looking for excuses to explore; I particularly enjoy San Francisco, Santa Cruz, and Carmel-by-the-Sea.</p>

			<h2>Projects</h2>

			<p>I do a lot of hobby programming. Most of it isn't that useful, but it is kind of fun. Most of my coding these days is done in <a href="https://golang.org">Go</a>, which is a pretty great language. I'm also proficient in Java, Python, C#, and a number of other languages.</p>

			<p>Here are some of my favorites:</p>

			<ul>
				<li><a href="https://github.com/BranLwyd/bNotify">bNotify</a>: A pure-Go daemon &amp; Android app which allows you to push notifications to your phone from any computer. It uses <a href="https://developers.google.com/cloud-messaging/">GCM</a> to push the notifications. The notification data is encrypted &amp; authenticated in transit. However, this is just a toy project, so you probably shouldn't trust it; at the very least, every server uses the same key so compromise of one server leads to compromise of the authenticity/confidentiality of messages from all servers. It does retry nicely, however. It's good for sending best-effort notifications for events from your various servers.</li>
				<li><a href="https://github.com/BranLwyd/rss-download">rss-download</a>: A pure-Go daemon &amp; CLI tool which watches a set of RSS feeds and downloads any new files (stored as links) to a specified directory. It normally polls slowly, but each RSS feed can be configured to poll more quickly during a specified interval.</li>
				<li><a href="https://github.com/BranLwyd/tumblr-sync">tumblr-sync</a>: A Java library which provides access to various Tumblr objects as native Java objects, and allows those objects to be persisted locally.</li>
				<li><a href="https://github.com/BranLwyd/bdcpu16">bdcpu16</a>: A Java implementation of the <a href="https://raw.githubusercontent.com/gatesphere/demi-16/master/docs/dcpu-specs/dcpu-1-7.txt">DCPU-16</a> virtual machine from Notch's cancelled 0x10c. Includes support for all hardware defined by the specification, an assembler, and a debugger.</li>
			</ul>
		</div>
	</div>
</body>
</html>`

	style = `
body {
  margin: 0;
  background-color: #f5f5f5;
  font-family: Helvetica, "Roboto", Tahoma, sans-serif;
}

.ribbon {
  width: 100%;
  height: 128px;
  background-color: #555;
  border-bottom: 1px solid black;
}

.content {
  padding: 80px 56px;
  font-size: 14px;
  background-color: #ffffff;
  color: #424242;
  margin: -96px auto 80px;
  max-width: 939px;
  width: calc(80% - 138px);
  border: 1px solid black;

  line-height: 24px;
}

h3 {
  font-weight: normal;
  margin: 48px 0px 24px;
  font-size: 34px;
  /* For a closer match of the original template: */
  font-family: "Roboto", Helvetica, Tahoma, sans-serif;
  line-height: 40px;
}

p {
  margin: 0 0 16px;
}

.content {
  padding: 0;
  min-width: 30em;
}

.inner-content {
  padding: 2em 4em;
}

.header {
  padding: 2em;
  background: #ddd;
  border-bottom: 2px solid gray;
}

.header h1 {
  font-family: "Droid Serif", serif;
  text-align: center;
  letter-spacing: 3px;
}`
)
