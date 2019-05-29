package main

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	defaultDBTimeFormat = "YYYYMMDDHH24"
	day                 = "YYYYMMDD"
	hour                = "YYYYMMDDHH24"
	minute              = "YYYYMMDDHH24MI"
	month               = "YYYYMM"

	// Template for generation of Input/Output throughput of a logical server
	lsThroughputQueryTemplate = `SELECT CASE
		WHEN a.time IS NOT NULL THEN a.time
		WHEN b.time IS NOT NULL THEN b.time
		ELSE NULL
		END as time,
			COALESCE(a.input_files, 0) AS input_files,
			COALESCE(b.input_cdrs, 0) AS input_cdrs,
			COALESCE(a.input_bytes, 0) AS input_bytes,
			COALESCE(b.output_files, 0) AS output_files,
			COALESCE(b.output_cdrs, 0) AS output_cdrs,
			COALESCE(b.output_bytes, 0) AS output_bytes
		FROM   (SELECT To_char(intime, '{{.TimeFormat}}') AS time,
			Count (*)                     AS input_files,
			Sum(bytes)::bigint               AS input_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '{{.TimeFormat}}') >= '{{.StartTime}}'
		AND to_char(intime, '{{.TimeFormat}}') <= '{{.EndTime}}'
		AND event = 67
		GROUP  BY To_char(intime, '{{.TimeFormat}}')
		ORDER  BY To_char(intime, '{{.TimeFormat}}')) a
		FULL OUTER JOIN (SELECT CASE
		WHEN c.time IS NOT NULL THEN c.time
		WHEN d.time IS NOT NULL THEN d.time
		ELSE NULL
		END AS time,
			c.input_cdrs,
			d.output_files,
			d.output_cdrs,
			d.output_bytes
		FROM   (SELECT To_char(intime, '{{.TimeFormat}}') AS time,
			COALESCE(Sum (cdrs)::bigint, 0)       AS
		input_cdrs
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '{{.TimeFormat}}') >= '{{.StartTime}}'
		AND to_char(intime, '{{.TimeFormat}}') <= '{{.EndTime}}'
		AND event = 73
		GROUP  BY To_char(intime, '{{.TimeFormat}}')
		ORDER  BY To_char(intime, '{{.TimeFormat}}')) c
		FULL OUTER JOIN (SELECT
		To_char(outtime, '{{.TimeFormat}}')
		AS time,
			Count(*)
		AS
		output_files,
			Sum(cdrs)::bigint
		AS
		output_cdrs,
			Sum(bytes)::bigint
		AS
		output_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(outtime, '{{.TimeFormat}}') >= '{{.StartTime}}'
		AND to_char(outtime, '{{.TimeFormat}}') <= '{{.EndTime}}'
		AND event = 68
		GROUP  BY To_char(outtime, '{{.TimeFormat}}')
		ORDER  BY To_char(outtime, '{{.TimeFormat}}')) d
		ON c.time = d.time) b
		ON a.time = b.time`

	// Template for generation of Input/Output throughput of a stream
	streamThroughputQueryTemplate = `SELECT CASE
		WHEN a.time IS NOT NULL THEN a.time
		WHEN b.time IS NOT NULL THEN b.time
		ELSE NULL
		END as time,
			COALESCE(a.total_input_files, 0) AS total_input_files,
			COALESCE(b.total_input_cdrs, 0) AS total_input_cdrs,
			COALESCE(a.total_input_bytes, 0) AS total_input_bytes,
			COALESCE(b.total_output_files, 0) AS total_output_files,
			COALESCE(b.total_output_cdrs, 0) AS total_output_cdrs,
			COALESCE(b.total_output_bytes, 0) AS total_output_bytes
		FROM   (SELECT To_char(intime, '{{.TimeFormat}}') AS time,
			Count (*)                     AS total_input_files,
			Sum(bytes)::bigint            AS total_input_bytes
		FROM   audittraillogentry
		WHERE  event = 67
			{{- $names := concat .InnodeNames -}}
			{{- $ids := concat .InnodeIds -}}
			{{- if and $names $ids -}}
				AND (trim(innodename) IN ({{- $names -}}) OR innodeid IN ({{- $ids -}}))
			{{- else if $names -}}
				AND (trim(innodename) IN ({{- $names -}}))
			{{- else if $ids -}}
				AND (innodeid IN ({{- $ids -}}))
			{{- else -}}
				AND 1=2
			{{- end -}}
		GROUP  BY To_char(intime, '{{.TimeFormat}}')
		ORDER  BY To_char(intime, '{{.TimeFormat}}')) a
		FULL OUTER JOIN (SELECT CASE
		WHEN c.time IS NOT NULL THEN c.time
		WHEN d.time IS NOT NULL THEN d.time
		ELSE NULL
		END AS time,
			c.total_input_cdrs,
			d.total_output_files,
			d.total_output_cdrs,
			d.total_output_bytes
		FROM   (SELECT To_char(intime, '{{.TimeFormat}}') AS time,
			COALESCE(Sum (cdrs)::bigint, 0)       AS
		total_input_cdrs
		FROM   audittraillogentry
		WHERE  event = 73
			{{- $names := concat .InnodeNames -}}
			{{- $ids := concat .InnodeIds -}}
			{{- if and $names $ids -}}
				AND (trim(innodename) IN ({{- $names -}}) OR innodeid IN ({{- $ids -}}))
			{{- else if $names -}}
				AND (trim(innodename) IN ({{- $names -}}))
			{{- else if $ids -}}
				AND (innodeid IN ({{- $ids -}}))
			{{- else -}}
				AND 1=2
			{{- end -}}
		GROUP  BY To_char(intime, '{{.TimeFormat}}')
		ORDER  BY To_char(intime, '{{.TimeFormat}}')) c
		FULL OUTER JOIN (SELECT
		To_char(outtime, '{{.TimeFormat}}')
		AS time,
			Count(*)
		AS
		total_output_files,
			Sum(cdrs)::bigint 
		AS
		total_output_cdrs,
			Sum(bytes)::bigint 
		AS
		total_output_bytes
		FROM   audittraillogentry
		WHERE  event = 68
			{{- $names := concat .OutnodeNames -}}
			{{- $ids := concat .OutnodeIds -}}
			{{- if and $names $ids -}}
				AND (trim(outnodename) IN ({{- $names -}}) OR outnodeid IN ({{- $ids -}}))
			{{- else if $names -}}
				AND (trim(outnodename) IN ({{- $names -}}))
			{{- else if $ids -}}
				AND (outnodeid IN ({{- $ids -}}))
			{{- else -}}
				AND 1=2
			{{- end -}}
		GROUP  BY To_char(outtime,
			'{{.TimeFormat}}'
		)
		ORDER  BY To_char(outtime,
			'{{.TimeFormat}}'
		)) d
		ON c.time = d.time) b
		ON a.time = b.time`
)

type AudittrailLogEntryQueryParameters struct {
	StartTime    string
	EndTime      string
	TimeFormat   string
	Collectors   []string
	Distributors []string
	InnodeNames  []string
	OutnodeNames []string
	InnodeIds    []string
	OutnodeIds   []string
}

func parseTemplate(templateName string, queryTemplate string, paramStruct interface{}) string {
	var actualQuery bytes.Buffer

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"concat": func(args []string) string {
			result := ""

			for i, arg := range args {
				result += fmt.Sprintf("'%s'", arg)

				if i+1 != len(args) {
					result += ","
				}
			}

			return result
		},
	}

	parsedTemplate := template.Must(template.New(templateName).Funcs(funcMap).Parse(queryTemplate))

	err := parsedTemplate.Execute(&actualQuery, paramStruct)

	if err != nil {
		fmt.Println(err)
	}

	return actualQuery.String()
}
