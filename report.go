package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/series"
	"github.com/jmoiron/sqlx"
	"github.com/kniren/gota/dataframe"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
	"strconv"
)

type ResultSet struct {
	columnsDataTypes map[string]series.Type
	data             dataframe.DataFrame
	table            *tablewriter.Table
}

func (r *ResultSet) GetColumnsNames() []string {
	return r.data.Names()
}

func (r *ResultSet) GetColumnsDataTypes() map[string]series.Type {
	return r.columnsDataTypes
}

func (r *ResultSet) GetColumnSeries(columnName string) series.Series {
	return r.data.Col(columnName)
}

func (r *ResultSet) WriteToTxtFile(filename string) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		logger.Error(err)
		return
	}

	table := tablewriter.NewWriter(file)
	table.SetHeader(r.data.Names())
	table.AppendBulk(r.data.Records()[1:])
	table.Render()
}

func (r *ResultSet) WriteToConsole() {
	if r.table == nil {
		r.table = tablewriter.NewWriter(os.Stdout)
	}

	r.table.SetHeader(r.data.Names())
	r.table.AppendBulk(r.data.Records()[1:])
	fmt.Fprintf(os.Stdout, "\n")
	r.table.Render()
}

func (r *ResultSet) WriteToCSVFile(filename string) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		logger.Error(err)
		return
	}

	w := csv.NewWriter(file)

	for _, record := range r.data.Records() {
		w.Write(record)
	}

	w.Flush()
}

type Report struct {
	name         string
	defaultTable *ResultSet
	avgTable     *ResultSet
	sumTable     *ResultSet
	minTable     *ResultSet
	maxTable     *ResultSet
}

func (r *Report) ExtractResultSet(rows *sqlx.Rows) {
	var row map[string]interface{}
	var rowFieldsStringVals []string
	var data [][]string

	columns, err := rows.Columns()

	if err != nil {
		fmt.Println(err)
	}

	if r.defaultTable == nil {
		r.defaultTable = &ResultSet{}
	}

	columnsTypes, err := rows.ColumnTypes()

	r.defaultTable.columnsDataTypes = map[string]series.Type{}

	for _, columnType := range columnsTypes {
		r.defaultTable.columnsDataTypes[columnType.Name()] = mapReflectTypeToSeriesType(columnType.ScanType())
	}

	data = append(data, columns)

	for rows.Next() {

		row = map[string]interface{}{}
		rowFieldsStringVals = []string{}

		rows.MapScan(row)

		for i := range columns {
			rowFieldsStringVals = append(rowFieldsStringVals, rowFieldToString(row[columns[i]]))
		}

		data = append(data, rowFieldsStringVals)
	}

	r.defaultTable.data = dataframe.LoadRecords(data, dataframe.WithTypes(r.defaultTable.columnsDataTypes))
}

func (r *Report) GetDefaultTable() *ResultSet {
	if r.defaultTable == nil {
		r.defaultTable = &ResultSet{}
	}

	return r.defaultTable
}

func (r *Report) GetAvgTable() *ResultSet {

	if r.avgTable == nil {
		var statsFields []string
		var records [][]string

		r.avgTable = &ResultSet{}

		records = append(records, r.GetDefaultTable().GetColumnsNames())

		for _, columnName := range r.GetDefaultTable().GetColumnsNames() {

			columnDataType := r.GetDefaultTable().GetColumnsDataTypes()[columnName]

			// Generate average for Float and Int columns only
			if columnDataType == series.Float || columnDataType == series.Int {
				if avg, err := stats.Mean(r.GetDefaultTable().GetColumnSeries(columnName).Float()); err == nil {
					statsFields = append(statsFields, strconv.FormatFloat(avg, 'f', -1, 64))
				}
			} else {
				statsFields = append(statsFields, "NA")
			}
		}

		records = append(records, statsFields)
		r.avgTable.data = dataframe.LoadRecords(records, dataframe.DefaultType(series.String), dataframe.DetectTypes(false))
	}

	return r.avgTable
}

func (r *Report) GetSumTable() *ResultSet {

	if r.avgTable == nil {
		var statsFields []string
		var records [][]string

		r.sumTable = &ResultSet{}
		records = append(records, r.GetDefaultTable().GetColumnsNames())

		for _, columnName := range r.defaultTable.GetColumnsNames() {
			// Generate average for Float and Int columns only
			if r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Float ||
				r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Int {

				if sum, err := stats.Sum(r.GetDefaultTable().GetColumnSeries(columnName).Float()); err == nil {
					statsFields = append(statsFields, strconv.FormatFloat(sum, 'f', -1, 64))
				}
			} else {
				statsFields = append(statsFields, "NA")
			}
		}
		records = append(records, statsFields)
		r.sumTable.data = dataframe.LoadRecords(records, dataframe.DefaultType(series.String), dataframe.DetectTypes(false))
	}

	return r.sumTable
}

func (r *Report) GetMinTable() *ResultSet {

	if r.minTable == nil {
		var statsFields []string
		var records [][]string

		r.minTable = &ResultSet{}
		records = append(records, r.GetDefaultTable().GetColumnsNames())

		for _, columnName := range r.defaultTable.GetColumnsNames() {
			// Generate average for Float and Int columns only
			if r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Float ||
				r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Int {

				if min, err := stats.Min(r.GetDefaultTable().GetColumnSeries(columnName).Float()); err == nil {
					statsFields = append(statsFields, strconv.FormatFloat(min, 'f', -1, 64))
				}
			} else {
				statsFields = append(statsFields, "NA")
			}
		}

		records = append(records, statsFields)
		r.minTable.data = dataframe.LoadRecords(records, dataframe.DefaultType(series.String), dataframe.DetectTypes(false))
	}

	return r.minTable
}

func (r *Report) GetMaxTable() *ResultSet {
	if r.maxTable == nil {
		var statsFields []string
		var records [][]string

		r.maxTable = &ResultSet{}
		records = append(records, r.GetDefaultTable().GetColumnsNames())

		for _, columnName := range r.GetDefaultTable().GetColumnsNames() {
			// Generate average for Float and Int columns only
			if r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Float ||
				r.GetDefaultTable().GetColumnsDataTypes()[columnName] == series.Int {

				if max, err := stats.Max(r.GetDefaultTable().GetColumnSeries(columnName).Float()); err == nil {
					statsFields = append(statsFields, strconv.FormatFloat(max, 'f', -1, 64))
				}
			} else {
				statsFields = append(statsFields, "NA")
			}
		}

		records = append(records, statsFields)
		r.maxTable.data = dataframe.LoadRecords(records, dataframe.DefaultType(series.String), dataframe.DetectTypes(false))
	}

	return r.maxTable
}

func mapReflectTypeToSeriesType(reflectType reflect.Type) series.Type {

	kind := reflectType.Kind()

	if kind == reflect.String {
		return series.String
	} else if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int32 || kind == reflect.Int64 ||
		kind == reflect.Uint || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return series.Int
	} else if kind == reflect.Bool {
		return series.Bool
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		return series.Float
	}

	return series.String
}

func rowFieldToString(field interface{}) string {
	var v = reflect.ValueOf(field)
	var t = reflect.TypeOf(field)
	var fieldStringValue string

	switch t.Kind() {
	case reflect.Int64:
		fieldStringValue = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint8:
		fieldStringValue = strconv.FormatInt(v.Int(), 10)
	case reflect.String:
		fieldStringValue = v.String()
	}

	return fieldStringValue
}

func rowFieldToFloat(field interface{}) (float64, error) {
	var v = reflect.ValueOf(field)
	var t = reflect.TypeOf(field)
	var fieldFloatValue float64

	switch t.Kind() {
	case reflect.Int64:
		fieldFloatValue = float64(v.Int())
	case reflect.Uint8:
		fieldFloatValue = float64(v.Int())
	default:
		return 0, fmt.Errorf("Not numeric field")
	}

	return fieldFloatValue, nil
}
