package database

import (
	config "emm-statistics-report/configuration"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

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

func GetStreamProcessedInOut(stream *config.Stream) *sqlx.Rows {

	if ls := config.FindLsRunningStream(stream); ls != nil {

		logger.WithFields(logrus.Fields{
			"logical_server": ls.Name,
			"stream":         stream.Name,
		}).Debug("Stream is active in a logical server")

		session := CreateSession(ls)

		collectorsInStatement := createSQLInStatement(stream.Collectors)
		distributorsInStatement := createSQLInStatement(stream.Distributors)

		query := fmt.Sprintf(`SELECT a.*, 
       	b.*, 
       	c.* 
		FROM   
		(SELECT Count (*)  AS total_input_files, 
               Sum(bytes) AS total_input_bytes 
        FROM   audittraillogentry 
        WHERE  event = 67 and trim(innodename) in %s) a,
         
       (SELECT COALESCE(Sum (cdrs), 0) AS total_input_cdrs 
        FROM   audittraillogentry 
        WHERE  event = 73 and trim(innodename) in %s) b, 

       (SELECT Count(*)  AS total_output_files, 
               Sum(cdrs) AS total_output_cdrs,
               Sum(bytes) AS total_output_bytes 
        FROM   audittraillogentry 
        WHERE  event = 68 and trim(outnodename) in %s) c;`,
			collectorsInStatement,
			collectorsInStatement,
			distributorsInStatement)

		logger.WithFields(logrus.Fields{
			"query": query,
		}).Debug("GetStreamProcessedInOut query")

		if session != nil {
			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": query,
					"error": err,
				}).Error("Stream processed input/output")
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

	for index, element := range inStatElements {
		inStatement += fmt.Sprintf("'%s'", element)

		if index < len(inStatElements)-1 {
			inStatement += ","
		}
	}

	return inStatement
}