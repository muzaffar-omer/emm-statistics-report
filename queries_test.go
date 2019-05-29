package main

import (
	"testing"
)

func TestParseTemplate(t *testing.T) {
	var query string
	var queryParams AudittrailLogEntryQueryParameters

	templateText := `
		{{- $names := concat .InnodeNames -}}
		{{- $ids := concat .InnodeIds -}}
		{{- if and $names $ids -}}
			AND (innodenames IN ({{- $names -}}) OR innodeids IN ({{- $ids -}}))
		{{- else if $names -}}
			AND (innodenames IN ({{- $names -}}))
		{{- else if $ids -}}
			AND (innodeids IN ({{- $ids -}}))
		{{- end -}}`

	// Empty Names and Ids
	queryParams = AudittrailLogEntryQueryParameters{
		InnodeNames: []string{},
		InnodeIds:   []string{},
	}
	query = parseTemplate("", templateText, queryParams)
	if query != "" {
		t.Errorf("Expecting '', but got '%s'", query)
	}

	// Names only
	queryParams = AudittrailLogEntryQueryParameters{
		InnodeNames: []string{"node1", "node2"},
	}
	query = parseTemplate("", templateText, queryParams)
	if query != "AND (innodenames IN (node1,node2))" {
		t.Errorf("Expecting 'AND (innodenames IN (node1,node2))', but got '%s'", query)
	}

	// Names and Ids
	queryParams = AudittrailLogEntryQueryParameters{
		InnodeNames: []string{"node1", "node2"},
		InnodeIds:   []string{"10", "20"},
	}
	query = parseTemplate("", templateText, queryParams)
	if query != "AND (innodenames IN (node1,node2) OR innodeids IN (10,20))" {
		t.Errorf("Expecting 'AND (innodenames IN (node1,node2) OR innodeids IN (10,20))', but got '%s'", query)
	}
}
