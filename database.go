package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"

	//"strconv"

	//"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var sessionsPool []Session

type Session struct {
	logicalServer *LogicalServer
	Db            *sqlx.DB
}

// Checks whether there is already an existing working session
// If it finds one, returns pointer to the session, otherwise, returns nil
func isSessionExists(logicalServer *LogicalServer) *Session {

	for _, existingSession := range sessionsPool {
		if existingSession.logicalServer.Equals(logicalServer) {
			return &existingSession
		}
	}

	return nil
}

// Opens a session to a logical server database and adds the
// session to the pool
func CreateSession(ls *LogicalServer) *Session {

	var newSession *Session

	if newSession = isSessionExists(ls); newSession != nil {
		return newSession
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=disable", ls.Username, ls.Database, ls.Password, ls.Port, ls.IP)

	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"logical_server": ls.Name,
			"database":       ls.Database,
			"ip":             ls.IP,
			"port":           ls.Port,
			"error":          err,
		}).Error("Opening session")

		return nil

	} else {
		newSession = &Session{logicalServer: ls, Db: db}
	}

	err = db.Ping()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"logical_server": ls.Name,
			"database":       ls.Database,
			"ip":             ls.IP,
			"port":           ls.Port,
			"error":          err,
		}).Error("Ping database")

		return nil
	}

	sessionsPool = append(sessionsPool, *newSession)

	return newSession
}

func (s Session) executeQuery(query string) *Report {
	var report Report

	rows, err := s.Db.Queryx(query)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"query": fmt.Sprintf(query),
			"error": err,
		}).Error("Querying all rows")
	}

	report = Report{}
	report.ExtractResultSet(rows)

	return &report
}

func printResultTable(rows *sqlx.Rows, caption string) [][]string {

	var table = tablewriter.NewWriter(os.Stdout)
	var statsMap = make(map[string][]float64)

	var row map[string]interface{}
	var rowFieldsStringVals []string
	var data [][]string

	columns, _ := rows.Columns()
	data = append(data, columns)

	for rows.Next() {
		row = make(map[string]interface{})
		rowFieldsStringVals = []string{}

		rows.MapScan(row)

		for i := range columns {
			rowFieldsStringVals = append(rowFieldsStringVals, rowFieldToString(row[columns[i]]))

			if floatFieldValue, err := rowFieldToFloat(row[columns[i]]); err == nil {
				statsMap[columns[i]] = append(statsMap[columns[i]], floatFieldValue)
			}
		}

		table.Append(rowFieldsStringVals)
		data = append(data, rowFieldsStringVals)
	}

	table.SetCaption(true, caption)
	table.SetHeader(columns)
	table.Render()

	// Average
	fmt.Printf("\n\nAverage\n")
	table = createAverageTable(columns, statsMap)
	table.Render()

	// Max
	fmt.Printf("\n\nMax\n")
	table = createMaxTable(columns, statsMap)
	table.Render()

	// Min
	fmt.Printf("\n\nMin\n")
	table = createMinTable(columns, statsMap)
	table.Render()

	// Sum
	fmt.Printf("\n\nSum\n")
	table = createSumTable(columns, statsMap)
	table.Render()

	return data
}

func createAverageTable(columns []string, statsMap map[string][]float64) *tablewriter.Table {
	var table = tablewriter.NewWriter(os.Stdout)
	var statsFields = []string{}

	table.SetHeader(columns)

	for i := range columns {
		if avg, err := stats.Mean(statsMap[columns[i]]); err == nil {
			statsFields = append(statsFields, strconv.FormatFloat(avg, 'f', -1, 64))
		} else {
			statsFields = append(statsFields, "")
		}
	}

	table.Append(statsFields)

	return table
}

func createMaxTable(columns []string, statsMap map[string][]float64) *tablewriter.Table {
	var table = tablewriter.NewWriter(os.Stdout)
	var statsFields = []string{}

	table.SetHeader(columns)

	for i := range columns {
		if avg, err := stats.Max(statsMap[columns[i]]); err == nil {
			statsFields = append(statsFields, strconv.FormatFloat(avg, 'f', -1, 64))
		} else {
			statsFields = append(statsFields, "")
		}
	}

	table.Append(statsFields)

	return table
}

func createMinTable(columns []string, statsMap map[string][]float64) *tablewriter.Table {
	var table = tablewriter.NewWriter(os.Stdout)
	var statsFields = []string{}

	table.SetHeader(columns)

	for i := range columns {
		if avg, err := stats.Min(statsMap[columns[i]]); err == nil {
			statsFields = append(statsFields, strconv.FormatFloat(avg, 'f', -1, 64))
		} else {
			statsFields = append(statsFields, "")
		}
	}

	table.Append(statsFields)

	return table
}

func createSumTable(columns []string, statsMap map[string][]float64) *tablewriter.Table {
	var table = tablewriter.NewWriter(os.Stdout)
	var statsFields = []string{}

	table.SetHeader(columns)

	for i := range columns {
		if avg, err := stats.Sum(statsMap[columns[i]]); err == nil {
			statsFields = append(statsFields, strconv.FormatFloat(avg, 'f', -1, 64))
		} else {
			statsFields = append(statsFields, "")
		}
	}

	table.Append(statsFields)

	return table
}

