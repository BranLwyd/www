package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

var (
	inFlag       = flag.String("in", "", "The input markdown file to convert.")
	outFlag      = flag.String("out", "", "The output HTML file to write.")
	templateFlag = flag.String("template", "", "The template file to place generated HTML in.")
)

func main() {
	// Parse & verify flags.
	flag.Parse()
	if *inFlag == "" {
		die("--in is required")
	}
	if *outFlag == "" {
		die("--out is required")
	}
	if *templateFlag == "" {
		die("--template is required")
	}

	// Read input & convert to HTML.
	md, err := ioutil.ReadFile(*inFlag)
	if err != nil {
		die("Couldn't read markdown %q: %v", *inFlag, err)
	}
	content := markdown.ToHTML(md, nil, html.NewRenderer(html.RendererOptions{Flags: html.HrefTargetBlank}))

	// Read template, replace content marker with content, & write result..
	tmpl, err := ioutil.ReadFile(*templateFlag)
	if err != nil {
		die("Couldn't read template %q: %v", *templateFlag, err)
	}
	html := bytes.Replace(tmpl, []byte("<!--CONTENT-->"), content, 1)
	if err := ioutil.WriteFile(*outFlag, html, 0640); err != nil {
		die("Couldn't write HTML %q: %v", *outFlag, err)
	}
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
