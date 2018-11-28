package database

import (
	"fmt"
	"strconv"
)

/*
	ATEvent ::= [APPLICATION 43] ENUMERATED
				{
        			collection(67),   -- 'C'
        			processing(80),   -- 'P'
        			inProcessing(73), -- 'I'
        			distribution(68), -- 'D'
        			any(65)           -- 'A'
				}
*/

type AuditTrailLogEntry struct {
	Event                    string
	Innodeid                 string
	Innodename               string
	Sourceid                 string
	Intime                   string
	Outnodeid                string
	Outnodename              string
	Destinationid            string
	Outtime                  string
	Bytes                    string
	Cdrs                     string
	Tableindex               string
	Noofsubfilesinfile       string
	Recordsequencenumberlist string
}

type TableDescription struct {
	TableSchema         string `db:"table_schema"`
	TableName           string `db:"table_name"`
	UserDefinedTypeName string `db:"user_defined_type_name"`
}

type TotalProcessedInOut struct {
	TotalInputFiles  int `db:"total_input_files"`
	TotalInputBytes  int `db:"total_input_bytes"`
	TotalInputCdrs   int `db:"total_input_cdrs"`
	TotalOutputFiles int `db:"total_output_files"`
	TotalOutputCdrs  int `db:"total_output_cdrs"`
	TotalOutputBytes int `db:"total_output_bytes"`
}

func (t TotalProcessedInOut) String() string {

	return fmt.Sprintf("total_input_files : %s, "+
		"total_input_bytes : %s, "+
		"total_input_cdrs: %s, "+
		"total_output_files : %s, "+
		"total_output_cdrs : %s, "+
		"total_output_bytes : %s\n",
		t.TotalInputFiles,
		t.TotalInputBytes,
		t.TotalInputCdrs,
		t.TotalOutputFiles,
		t.TotalOutputCdrs,
		t.TotalOutputBytes)
}

func (t TotalProcessedInOut) AsArray() []int {
	return []int{t.TotalInputFiles, t.TotalOutputBytes, t.TotalOutputCdrs, t.TotalInputFiles, t.TotalOutputCdrs,
		t.TotalOutputBytes}
}

func (t TotalProcessedInOut) Header() []string {
	return []string{"total_input_files", "total_input_bytes", "total_input_cdrs", "total_output_files", "total_output_cdrs",
		"total_output_bytes"}
}

type TotalGroupedProcessedInOut struct {
	Time             string `db:"time"`
	TotalInputFiles  int    `db:"total_input_files"`
	TotalInputCdrs   int    `db:"total_input_cdrs"`
	TotalInputBytes  int    `db:"total_input_bytes"`
	TotalOutputFiles int    `db:"total_output_files"`
	TotalOutputCdrs  int    `db:"total_output_cdrs"`
	TotalOutputBytes int    `db:"total_output_bytes"`
}

func (t TotalGroupedProcessedInOut) String() string {

	return fmt.Sprintf("time : %s,"+
		"total_input_files : %s, "+
		"total_input_cdrs: %s, "+
		"total_input_bytes : %s, "+
		"total_output_files : %s, "+
		"total_output_cdrs : %s, "+
		"total_output_bytes : %s\n",
		t.Time,
		t.TotalInputFiles,
		t.TotalInputBytes,
		t.TotalInputCdrs,
		t.TotalOutputFiles,
		t.TotalOutputCdrs,
		t.TotalOutputBytes)
}

func (t TotalGroupedProcessedInOut) AsArray() []string {
	return []string{t.Time,
		strconv.Itoa(t.TotalInputFiles),
		strconv.Itoa(t.TotalInputCdrs),
		strconv.Itoa(t.TotalInputBytes),
		strconv.Itoa(t.TotalOutputFiles),
		strconv.Itoa(t.TotalOutputCdrs),
		strconv.Itoa(t.TotalOutputBytes)}
}

func (t TotalGroupedProcessedInOut) Header() []string {
	return []string{"time", "total_input_files", "total_input_cdrs", "total_input_bytes", "total_output_files", "total_output_cdrs",
		"total_output_bytes"}
}

func (t TotalGroupedProcessedInOut) GetStatisticsMap() map[string]int {
	statsMap := make(map[string]int)

	statsMap["Input Files"] = t.TotalInputFiles
	statsMap["Input CDRs"] = t.TotalInputCdrs
	statsMap["Input Bytes"] = t.TotalInputBytes
	statsMap["Output Files"] = t.TotalOutputFiles
	statsMap["Output CDRs"] = t.TotalOutputCdrs
	statsMap["Output Bytes"] = t.TotalOutputBytes

	return statsMap
}
