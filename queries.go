package main

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	defaultDBTimeFormat = "YYYYMMDDHH24"
	day = "YYYYMMDD"
	hour = "YYYYMMDDHH24"
	minute = "YYYYMMDDHH24MI"
	month = "YYYYMM"

	// Template for generation of Input/Output throughput of a logical server
	throughputQueryTemplate = `SELECT CASE
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

type TrafficQueryParameters struct {
}

func parseTemplate(templateName string, queryTemplate string, paramStruct interface{}) string {
	var actualQuery bytes.Buffer
	var parsedTemplate = template.Must(template.New(templateName).Parse(queryTemplate))

	err := parsedTemplate.Execute(&actualQuery, paramStruct)

	if err != nil {
		fmt.Println(err)
	}

	return actualQuery.String()
}