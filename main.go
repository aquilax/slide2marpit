package main

import (
	"bufio"
	"html/template"
	"log"
	"os"

	"golang.org/x/tools/present"
)

var mdTmpl = `
{{define "root" -}}
---

---

# {{.Title}}
{{with .Subtitle}}## {{.}}</h2>
{{end -}}
{{range .Authors}}<author>
{{range .Elem}}{{elem $.Template .}}{{end}}</author>
{{end -}}

---

{{range .Sections}}{{elem $.Template .}}
{{end -}}
{{end}}

{{define "newline"}}{{/* No automatic line break. Paragraphs are free-form. */}}
{{end}}

{{define "section"}}
{{if .Title}}## {{.Title}}

{{end -}}
{{range .Elem}}{{elem $.Template .}}{{end}}
---{{- end}}

{{define "list" -}}
{{range .Bullet -}}
* {{.}}
{{end -}}
{{end}}

{{define "text" -}}
{{if .Pre -}}
` + "```" + `{{range .Lines}}{{.}}{{end}}` + "```" + `
{{else -}}
{{range $i, $l := .Lines}}{{if $i}}{{template "newline"}}{{end}}{{style $l}}{{end}}
{{end -}}
{{end}}

{{define "code" -}}
{{if .Play -}}
<div class="playground">{{.Text}}</div>
{{else -}}
` + "```go" + `{{.Text}}` + "```" + `
{{end -}}
{{end}}

{{define "image" -}}
![{{with .Height}}height:{{.}}px{{end}} {{with .Width}}width:{{.}}px{{end}}]({{.URL}})
{{end}}

{{define "caption" -}}
<figcaption>{{style .Text}}</figcaption>
{{end}}

{{define "iframe" -}}
<iframe src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}></iframe>
{{end}}

{{define "link" -}}
[{{.Label}}]({{.URL}})
{{end}}

{{define "html" -}}{{.HTML}}{{end}}
`

func main() {
	tmpl := template.Must(present.Template().Parse(mdTmpl))
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	doc, err := present.Parse(reader, "name", 0)
	if err != nil {
		log.Fatal(err)
	}
	doc.Render(os.Stdout, tmpl)
}
