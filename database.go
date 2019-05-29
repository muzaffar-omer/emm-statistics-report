package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
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

func executeQuery(ls *LogicalServer, query string) *sqlx.Rows {

	var rows *sqlx.Rows

	if session := CreateSession(ls); session != nil {

		rows, err := session.Db.Queryx(query)

		if err != nil {
			logger.WithFields(logrus.Fields{
				"query": fmt.Sprintf(query),
				"error": err,
			}).Error("Querying all rows")
		}

		return rows

	} else {
		fmt.Println("Session is nil")
	}

	return rows
}

func printResultTable(rows *sqlx.Rows, caption string) {

	var table = tablewriter.NewWriter(os.Stdout)
	var statsMap = make(map[string][]float64)

	var row map[string]interface{}
	var rowFieldsStringVals []string

	columns, _ := rows.Columns()

	table.SetCaption(true, caption)
	table.SetHeader(columns)

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
	}

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

func createMaxTable(columns []string,  statsMap map[string][]float64) *tablewriter.Table {
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

func createMinTable(columns []string,  statsMap map[string][]float64) *tablewriter.Table {
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

func createSumTable(columns []string,  statsMap map[string][]float64) *tablewriter.Table {
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
