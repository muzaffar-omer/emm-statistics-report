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
			distributorsInStatement,
			distributorsInStatement)

		if session != nil {
			rows, err := session.Db.Queryx(query)

			if err != nil {
				logger.WithFields(logrus.Fields{
					"query": query,
				}).Error("Stream processed input/output")
			} else {
				return rows
			}
		}

	} else {
		logger.WithFields(logrus.Fields{
			"stream": stream.Name,
		}).Error("No logical server running the stream")
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

	inStatement := " ("

	for index, element := range inStatElements {
		inStatement += fmt.Sprintf("'%s'", element)

		if index < len(inStatElements)-1 {
			inStatement += ","
		}
	}

	inStatement += ") "

	return inStatement
}
