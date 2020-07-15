package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	inFlag       = flag.String("in", "", "The input markdown file to convert.")
	outFlag      = flag.String("out", "", "The output HTML file to write.")
	templateFlag = flag.String("template", "", "The template file to place generated HTML in.")
	titleFlag    = flag.String("title", "", "The title of the generated HTML page.")

	titleMarker   = []byte("<!--TITLE-->")
	contentMarker = []byte("<!--CONTENT-->")
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
	if *titleFlag == "" {
		die("--title is required")
	}

	// Read input & convert to HTML.
	md, err := ioutil.ReadFile(*inFlag)
	if err != nil {
		die("Couldn't read markdown %q: %v", *inFlag, err)
	}
	content := markdown.ToHTML(md, parser.NewWithExtensions(parser.CommonExtensions & ^parser.MathJax), html.NewRenderer(html.RendererOptions{Flags: html.HrefTargetBlank}))

	// Read template, replace title & content markers with content, and write result.
	htmlTmplBytes, err := ioutil.ReadFile(*templateFlag)
	if err != nil {
		die("Couldn't read template %q: %v", *templateFlag, err)
	}
	htmlTmpl, err := template.New("").Parse(string(htmlTmplBytes))
	if err != nil {
		die("Couldn't parse template: %v", err)
	}
	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, struct {
		Title   string
		Content string
	}{
		Title:   *titleFlag,
		Content: string(content),
	}); err != nil {
		die("Couldn't execute template: %v", err)
	}
	if err := ioutil.WriteFile(*outFlag, htmlBuf.Bytes(), 0640); err != nil {
		die("Couldn't write HTML %q: %v", *outFlag, err)
	}
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
