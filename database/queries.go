package database

import (
	config "emm-statistics-report/configuration"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

const DB_DATE_FORMAT = "YYYY-MM-DD"

var logger = config.Log()

func ListTables(session *Session) *sqlx.Rows {

	query := fmt.Sprintf("select table_schema, table_name, user_defined_type_name from information_schema.tables")

	rows, err := session.Db.Queryx(query)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"query": query,
			"error": err,
		}).Error("Describe tables")

		return nil
	}

	return rows
}

func GetStreamTotalProcessedInOut(logicalServer *config.LogicalServer, stream *config.Stream) *sqlx.Rows {

	if logicalServer != nil {
		//fmt.Println(logicalServer.Name)

		if session := CreateSession(logicalServer); session != nil {
			query := `SELECT a.*, 
       	b.*, 
       	c.* 
		FROM   
		(SELECT Count (*)  AS total_input_files, 
               Sum(bytes) AS total_input_bytes 
        FROM   audittraillogentry 
        WHERE  event = 67) a,
         
       (SELECT COALESCE(Sum (cdrs), 0) AS total_input_cdrs 
        FROM   audittraillogentry 
        WHERE  event = 73) b, 

       (SELECT Count(*)  AS total_output_files, 
               Sum(cdrs) AS total_output_cdrs,
               Sum(bytes) AS total_output_bytes 
        FROM   audittraillogentry 
        WHERE  event = 68) c;`

			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": fmt.Sprintf(query),
					"error": err,
				}).Error("Querying all rows")

				return nil
			}

			return rows
		}
	}

	return nil
}

func GetTotalProcessedInOut(session *Session) *sqlx.Rows {
	query := `SELECT a.*, 
       	b.*, 
       	c.* 
		FROM   
		(SELECT Count (*)  AS total_input_files, 
               Sum(bytes) AS total_input_bytes 
        FROM   audittraillogentry 
        WHERE  event = 67) a,
         
       (SELECT COALESCE(Sum (cdrs), 0) AS total_input_cdrs 
        FROM   audittraillogentry 
        WHERE  event = 73) b, 

       (SELECT Count(*)  AS total_output_files, 
               Sum(cdrs) AS total_output_cdrs,
               Sum(bytes) AS total_output_bytes 
        FROM   audittraillogentry 
        WHERE  event = 68) c;`

	rows, err := session.Db.Queryx(query)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"query": fmt.Sprintf(query),
			"error": err,
		}).Error("Querying all rows")

		return nil
	}

	return rows
}

func GetGroupedStreamProcessedInOut(stream *config.Stream, groupBy string) *sqlx.Rows {

	if ls := config.FindLsRunningStream(stream); ls != nil {

		logger.WithFields(logrus.Fields{
			"logical_server": ls.Name,
			"stream":         stream.Name,
		}).Debug("Stream is active in a logical server")

		session := CreateSession(ls)

		collectorsInStatement := createSQLInStatement(stream.Collectors)
		distributorsInStatement := createSQLInStatement(stream.Distributors)

		// By default, group by day
		dateFormat := "YYYY-MM-DD"

		switch groupBy {
		case "minute":
			dateFormat = "YYYY-MM-DD HH24MI"
		case "hour":
			dateFormat = "YYYY-MM-DD HH24"
			break
		case "day":
			dateFormat = "YYYY-MM-DD"
			break
		case "month":
			dateFormat = "YYYY-MM"
			break
		}

		query := fmt.Sprintf(`SELECT CASE
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
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			Count (*)                     AS total_input_files,
			Sum(bytes)                    AS total_input_bytes
		FROM   audittraillogentry
		WHERE  event = 67
		AND Trim(innodename) IN ( %[2]s )
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) a
		FULL OUTER JOIN (SELECT CASE
		WHEN c.time IS NOT NULL THEN c.time
		WHEN d.time IS NOT NULL THEN d.time
		ELSE NULL
		END AS time,
			c.total_input_cdrs,
			d.total_output_files,
			d.total_output_cdrs,
			d.total_output_bytes
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			COALESCE(Sum (cdrs), 0)       AS
		total_input_cdrs
		FROM   audittraillogentry
		WHERE  event = 73
		AND Trim(innodename) IN (
			%[2]s )
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) c
		FULL OUTER JOIN (SELECT
		To_char(outtime, '%[1]s')
		AS time,
			Count(*)
		AS
		total_output_files,
			Sum(cdrs)
		AS
		total_output_cdrs,
			Sum(bytes)
		AS
		total_output_bytes
		FROM   audittraillogentry
		WHERE  event = 68
		AND Trim(outnodename) IN
		(
			%[3]s )
		GROUP  BY To_char(outtime,
			'%[1]s'
		)
		ORDER  BY To_char(outtime,
			'%[1]s'
		)) d
		ON c.time = d.time) b
		ON a.time = b.time`, dateFormat, collectorsInStatement, distributorsInStatement)

		logger.WithFields(logrus.Fields{
			"query": query,
		}).Debug("GetGroupedStreamProcessedInOut query")

		if session != nil {
			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": query,
					"error": err,
				}).Error("Stream grouped processed input/output")
			} else {

				return rows
			}
		}

	} else {
		logger.WithFields(logrus.Fields{
			"stream": stream.Name,
		}).Error("Stream not running in any logical server")
	}

	return nil
}

func GetStreamProcessedInOut(stream *config.Stream, groupBy string, startDate time.Time, endDate time.Time) *sqlx.Rows {

	if ls := config.FindLsRunningStream(stream); ls != nil {

		logger.WithFields(logrus.Fields{
			"logical_server": ls.Name,
			"stream":         stream.Name,
		}).Debug("Stream is active in a logical server")

		session := CreateSession(ls)

		collectorsInStatement := createSQLInStatement(stream.Collectors)
		distributorsInStatement := createSQLInStatement(stream.Distributors)

		// By default, group by day
		dateFormat := "YYYY-MM-DD"

		switch groupBy {
		case "minute":
			dateFormat = "YYYY-MM-DD HH24MI"
		case "hour":
			dateFormat = "YYYY-MM-DD HH24"
			break
		case "day":
			dateFormat = "YYYY-MM-DD"
			break
		case "month":
			dateFormat = "YYYY-MM"
			break
		}

		query := fmt.Sprintf(`SELECT CASE
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
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			Count (*)                     AS total_input_files,
			Sum(bytes)                    AS total_input_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '%[4]s') >= '%[5]s'
		AND to_char(intime, '%[4]s') <= '%[6]s'
		AND event = 67
		AND Trim(innodename) IN ( %[2]s )
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) a
		FULL OUTER JOIN (SELECT CASE
		WHEN c.time IS NOT NULL THEN c.time
		WHEN d.time IS NOT NULL THEN d.time
		ELSE NULL
		END AS time,
			c.total_input_cdrs,
			d.total_output_files,
			d.total_output_cdrs,
			d.total_output_bytes
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			COALESCE(Sum (cdrs), 0)       AS
		total_input_cdrs
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '%[4]s') >= '%[5]s'
		AND to_char(intime, '%[4]s') <= '%[6]s'
		AND event = 73
		AND Trim(innodename) IN (
			%[2]s )
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) c
		FULL OUTER JOIN (SELECT
		To_char(outtime, '%[1]s')
		AS time,
			Count(*)
		AS
		total_output_files,
			Sum(cdrs)
		AS
		total_output_cdrs,
			Sum(bytes)
		AS
		total_output_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(outtime, '%[4]s') >= '%[5]s'
		AND to_char(outtime, '%[4]s') <= '%[6]s'
		AND event = 68
		AND Trim(outnodename) IN
		(
			%[3]s )
		GROUP  BY To_char(outtime,
			'%[1]s'
		)
		ORDER  BY To_char(outtime,
			'%[1]s'
		)) d
		ON c.time = d.time) b
		ON a.time = b.time`, dateFormat,
			collectorsInStatement,
			distributorsInStatement,
			DB_DATE_FORMAT,
			ConvertTimeToDBFormat(startDate),
			ConvertTimeToDBFormat(endDate))

		logger.WithFields(logrus.Fields{
			"query": query,
		}).Debug("GetGroupedStreamProcessedInOut query")

		if session != nil {
			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": query,
					"error": err,
				}).Error("Stream grouped processed input/output")
			} else {

				return rows
			}
		}

	} else {
		logger.WithFields(logrus.Fields{
			"stream": stream.Name,
		}).Error("Stream not running in any logical server")
	}

	return nil
}

func GetLogicalServerProcessedInOut(ls *config.LogicalServer, groupBy string, startDate time.Time, endDate time.Time) *sqlx.Rows {

	if ls != nil {

		session := CreateSession(ls)

		// By default, group by day
		dateFormat := "YYYY-MM-DD"

		switch groupBy {
		case "minute":
			dateFormat = "YYYY-MM-DD HH24MI"
		case "hour":
			dateFormat = "YYYY-MM-DD HH24"
			break
		case "day":
			dateFormat = "YYYY-MM-DD"
			break
		case "month":
			dateFormat = "YYYY-MM"
			break
		}

		query := fmt.Sprintf(`SELECT CASE
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
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			Count (*)                     AS total_input_files,
			Sum(bytes)                    AS total_input_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '%[2]s') >= '%[3]s'
		AND to_char(intime, '%[2]s') <= '%[4]s'
		AND event = 67
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) a
		FULL OUTER JOIN (SELECT CASE
		WHEN c.time IS NOT NULL THEN c.time
		WHEN d.time IS NOT NULL THEN d.time
		ELSE NULL
		END AS time,
			c.total_input_cdrs,
			d.total_output_files,
			d.total_output_cdrs,
			d.total_output_bytes
		FROM   (SELECT To_char(intime, '%[1]s') AS time,
			COALESCE(Sum (cdrs), 0)       AS
		total_input_cdrs
		FROM   audittraillogentry
		WHERE  
		to_char(intime, '%[2]s') >= '%[3]s'
		AND to_char(intime, '%[2]s') <= '%[4]s'
		AND event = 73
		GROUP  BY To_char(intime, '%[1]s')
		ORDER  BY To_char(intime, '%[1]s')) c
		FULL OUTER JOIN (SELECT
		To_char(outtime, '%[1]s')
		AS time,
			Count(*)
		AS
		total_output_files,
			Sum(cdrs)
		AS
		total_output_cdrs,
			Sum(bytes)
		AS
		total_output_bytes
		FROM   audittraillogentry
		WHERE  
		to_char(outtime, '%[2]s') >= '%[3]s'
		AND to_char(outtime, '%[2]s') <= '%[4]s'
		AND event = 68
		GROUP  BY To_char(outtime,
			'%[1]s'
		)
		ORDER  BY To_char(outtime,
			'%[1]s'
		)) d
		ON c.time = d.time) b
		ON a.time = b.time`, dateFormat,
			DB_DATE_FORMAT,
			ConvertTimeToDBFormat(startDate),
			ConvertTimeToDBFormat(endDate))

		logger.WithFields(logrus.Fields{
			"query": query,
		}).Debug("GetGroupedStreamProcessedInOut query")

		if session != nil {
			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": query,
					"error": err,
				}).Error("Stream grouped processed input/output")
			} else {

				return rows
			}
		}

	} else {
		logger.Error("Nil logical server")
	}

	return nil
}

func GetAllRows(session *Session, table string) *sqlx.Rows {

	rows, err := session.Db.Queryx(fmt.Sprintf("SELECT * FROM %s", table))

	if err != nil {
		logger.WithFields(logrus.Fields{
			"table": table,
			"query": fmt.Sprintf(`SELECT * FROM %s`, table),
			"error": err,
		}).Error("Querying all rows")

		return nil
	}

	return rows
}

func createSQLInStatement(inStatElements []string) string {

	inStatement := ""

	// In case there are no collectors or distributors, return empty quoted string
	if len(inStatElements) == 0 {
		inStatement = "''"
	}

	for index, element := range inStatElements {
		inStatement += fmt.Sprintf("'%s'", element)

		if index < len(inStatElements)-1 {
			inStatement += ","
		}
	}

	return inStatement
}

func ConvertCmdDateToTime(cmdDate string) (time.Time, error) {
	date, err := time.Parse("20060102", cmdDate)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"date":  cmdDate,
			"error": err,
		}).Error("Could not convert date")
	}

	return date, err
}

// Return date in the format YYYY-MM-DD
func ConvertTimeToDBFormat(time time.Time) string {
	return time.Format("2006-01-02")
}
