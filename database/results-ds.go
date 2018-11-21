package database

import (
	"fmt"
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
	TotalInputFiles  string `db:"total_input_files"`
	TotalInputBytes  string `db:"total_input_bytes"`
	TotalInputCdrs   string `db:"total_input_cdrs"`
	TotalOutputFiles string `db:"total_output_files"`
	TotalOutputCdrs  string `db:"total_output_cdrs"`
	TotalOutputBytes string `db:"total_output_bytes"`
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

func (t TotalProcessedInOut) AsArray() []string {
	return []string{t.TotalInputFiles, t.TotalOutputBytes, t.TotalOutputCdrs, t.TotalInputFiles, t.TotalOutputCdrs,
		t.TotalOutputBytes}
}

func (t TotalProcessedInOut) Header() []string {
	return []string{"total_input_files", "total_input_bytes", "total_input_cdrs", "total_output_files", "total_output_cdrs",
		"total_output_bytes"}
}

type TotalGroupedProcessedInOut struct {
	Time             string `db:"time"`
	TotalInputFiles  string `db:"total_input_files"`
	TotalInputBytes  string `db:"total_input_bytes"`
	TotalInputCdrs   string `db:"total_input_cdrs"`
	TotalOutputFiles string `db:"total_output_files"`
	TotalOutputCdrs  string `db:"total_output_cdrs"`
	TotalOutputBytes string `db:"total_output_bytes"`
}

func (t TotalGroupedProcessedInOut) String() string {

	return fmt.Sprintf("time : %s,"+
		"total_input_files : %s, "+
		"total_input_bytes : %s, "+
		"total_input_cdrs: %s, "+
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
	return []string{t.Time, t.TotalInputFiles, t.TotalOutputBytes, t.TotalOutputCdrs, t.TotalInputFiles, t.TotalOutputCdrs,
		t.TotalOutputBytes}
}

func (t TotalGroupedProcessedInOut) Header() []string {
	return []string{"time", "total_input_files", "total_input_bytes", "total_input_cdrs", "total_output_files", "total_output_cdrs",
		"total_output_bytes"}
}
