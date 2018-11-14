package database

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func ListTables(session Session) (*sql.Rows, error) {
	return nil, nil
}

func GetAllRows(session *Session, table string) *sql.Rows {

	rows, err := session.Db.Query(fmt.Sprintf("SELECT * FROM %s", table))

	if err != nil {
		log.WithFields(log.Fields{
			"table": table,
			"query": fmt.Sprintf(`SELECT * FROM %s`, table),
			"error": err,
		}).Error("Querying all rows from table")

		return nil
	}

	return rows
}
