package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
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
