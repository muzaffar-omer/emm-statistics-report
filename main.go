package main

import (
	"emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true})

	config := configuration.Config{}
	config.Init()

	log.WithFields(log.Fields{
		"ip":              config.Ip(),
		"port":            config.Port(),
		"password":        config.Password(),
		"username":        config.Username(),
		"logical_servers": config.Clusters(),
		"streams":         config.Streams(),
	}).Info("Configuration")

	//ls := configuration.LogicalServer{Ip: "pellefant.db.elephantsql.com", Port: "5432", Database:"opvvltfu",
	//Username: "opvvltfu", Password: "uacFLUmjHysihAvkd8lYizg0P9PcbD77"}

	ls := configuration.LogicalServer{Ip: "10.135.5.81", Port: "5432", Database: "fm_db_Server1", Username: "mmsuper", Password: "mediation"}

	session := database.CreateSession(ls)

	if session != nil {

		log.Info("Database opened successfully")
		//atlRow := database.AuditTrailLogEntry{}
		//tableDescription := database.TableDescription{}
		totalProcessedInOut := database.TotalProcessedInOut{}

		// rows := database.GetAllRows(session, "audittraillogentry")
		// rows := database.ListTables(session)
		rows := database.GetTotalProcessedInOut(session)
		noRows := 0

		//for rows.NextResultSet() {
		for rows.Next() {

			rows.StructScan(&totalProcessedInOut)
			// fmt.Println(tableDescription.TableName, tableDescription.TableSchema, tableDescription.UserDefinedTypeName)
			fmt.Println(totalProcessedInOut.TotalInputBytes, totalProcessedInOut.TotalInputCdrs, totalProcessedInOut.TotalInputFiles)

			noRows += 1
		}
		//}

		rows.Close()

		fmt.Printf("Total number of rows %d\n", noRows)
	}
}
