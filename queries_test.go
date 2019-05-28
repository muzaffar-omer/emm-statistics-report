package main

import (
	"bytes"
	"testing"
	"text/template"
)


func TestParseTemplate(t *testing.T) {
	var actualQuery bytes.Buffer

	type Employee struct {
		Username string
		Password string
		IDs []string
	}

	e := Employee{Username:"Muzaffar", Password:"Test", IDs: []string{"10", "20", "30"}}

	templateText := `
	{{- .Username}} && {{.Password}}
	{{if len .IDs gt 0}}
		Test
	{{end}}
	`

	parsedTemplate := template.Must(template.New("").Parse(templateText))

	parsedTemplate.Execute(&actualQuery, e)

	t.Log(actualQuery.String())
	t.Log("New Test")
}