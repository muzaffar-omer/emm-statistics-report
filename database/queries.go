package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func ListTables(session *Session) *sqlx.Rows {

	query := fmt.Sprintf("select table_schema, table_name, user_defined_type_name from information_schema.tables")

	rows, err := session.Db.Queryx(query)

	if err != nil {
		log.WithFields(log.Fields{
			"query": query,
			"error": err,
		}).Error("Describe tables")

		return nil
	}

	return rows
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
               Sum(cdrs) AS total_output_cdrs 
        FROM   audittraillogentry 
        WHERE  event = 68) c;`

	rows, err := session.Db.Queryx(query)

	if err != nil {
		log.WithFields(log.Fields{
			"query": fmt.Sprintf(query),
			"error": err,
		}).Error("Querying all rows")

		return nil
	}

	return rows
}

func GetAllRows(session *Session, table string) *sqlx.Rows {

	rows, err := session.Db.Queryx(fmt.Sprintf("SELECT * FROM %s", table))

	if err != nil {
		log.WithFields(log.Fields{
			"table": table,
			"query": fmt.Sprintf(`SELECT * FROM %s`, table),
			"error": err,
		}).Error("Querying all rows")

		return nil
	}

	return rows
}
