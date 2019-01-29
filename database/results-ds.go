package database

import (
	"fmt"
	"regexp"
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

	return fmt.Sprintf("total_input_files : %d, "+
		"total_input_bytes : %d, "+
		"total_input_cdrs: %d, "+
		"total_output_files : %d, "+
		"total_output_cdrs : %d, "+
		"total_output_bytes : %d\n",
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
		"total_input_files : %d, "+
		"total_input_cdrs: %d, "+
		"total_input_bytes : %d, "+
		"total_output_files : %d, "+
		"total_output_cdrs : %d, "+
		"total_output_bytes : %d\n",
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
		GetFormattedNumber(convertToMB(t.TotalInputBytes)),
		strconv.Itoa(t.TotalOutputFiles),
		strconv.Itoa(t.TotalOutputCdrs),
		GetFormattedNumber(convertToMB(t.TotalOutputBytes))}
}

func (t TotalGroupedProcessedInOut) Header() []string {
	return []string{"time", "total_input_files", "total_input_cdrs", "total_input_bytes_mb", "total_output_files", "total_output_cdrs",
		"total_output_bytes_mb"}
}

func (t TotalGroupedProcessedInOut) GetStatisticsMap() map[string]float64 {
	statsMap := make(map[string]float64)

	statsMap["Input Files"] = float64(t.TotalInputFiles)
	statsMap["Input CDRs"] = float64(t.TotalInputCdrs)
	statsMap["Input Bytes (MB)"] = convertToMB(t.TotalInputBytes)
	statsMap["Output Files"] = float64(t.TotalOutputFiles)
	statsMap["Output CDRs"] = float64(t.TotalOutputCdrs)
	statsMap["Output Bytes (MB)"] = convertToMB(t.TotalOutputBytes)

	return statsMap
}

func convertToMB(bytes int) float64 {
	return float64(bytes) / (1024 * 1024)
}

func GetFormattedNumber(number float64) string {
	if IsActualFloat(number) {
		return fmt.Sprintf("%f", number)
	} else {
		return fmt.Sprintf("%.0f", number)
	}
}

func IsActualFloat(number float64) bool {
	floatFormat := regexp.MustCompile("\\.([0-9]+)$")

	precisions := floatFormat.FindStringSubmatch(fmt.Sprintf("%f", number))

	if len(precisions) > 0 {

		precisionVal, err := strconv.Atoi(precisions[len(precisions)-1])

		if err != nil {
			fmt.Printf("%s\n", err)
		} else if precisionVal > 0 {
			return true
		}
	}

	return false
}
