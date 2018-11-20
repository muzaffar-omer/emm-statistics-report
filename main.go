package main

import (
	config "emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := config.Log()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true})

	logger.WithFields(logrus.Fields{
		"ip":       config.CmdConfig.Ip(),
		"port":     config.CmdConfig.Port(),
		"password": config.CmdConfig.Password(),
		"username": config.CmdConfig.Username(),
		"clusters": config.FileConfig.Clusters,
		"streams":  config.FileConfig.Streams,
	}).Info("Configuration")

	//ls := configuration.LogicalServer{Ip: "pellefant.db.elephantsql.com", Port: "5432", Database:"opvvltfu",
	//Username: "opvvltfu", Password: "uacFLUmjHysihAvkd8lYizg0P9PcbD77"}

	//ls := configuration.LogicalServer{Ip: "10.135.5.81", Port: "5781", Database: "fm_db_Server5", Username: "mmsuper", Password: "mediation"}
	//
	//session := database.CreateSession(&ls)
	//
	//if session != nil {
	//
	//	log.Info("Database opened successfully")
	//	//atlRow := database.AuditTrailLogEntry{}
	//	//tableDescription := database.TableDescription{}
	//	totalProcessedInOut := database.TotalProcessedInOut{}
	//
	//	// rows := database.GetAllRows(session, "audittraillogentry")
	//	// rows := database.ListTables(session)
	//	rows := database.GetTotalProcessedInOut(session)
	//	noRows := 0
	//
	//	//for rows.NextResultSet() {
	//	for rows.Next() {
	//
	//		rows.StructScan(&totalProcessedInOut)
	//		// fmt.Println(tableDescription.TableName, tableDescription.TableSchema, tableDescription.UserDefinedTypeName)
	//		fmt.Println(totalProcessedInOut)
	//
	//		noRows += 1
	//	}
	//
	//	rows.Close()
	//
	//	fmt.Printf("Total number of rows %d\n", noRows)
	//}
	//
	//stream := config.FileConfig.GetStreamConfig("SMSC_IC")
	//foundLs := config.FileConfig.FindLsRunningStream(stream)
	//
	//if foundLs != nil {
	//	fmt.Println(foundLs.Name)
	//}

	var totalProcessedInOut database.TotalProcessedInOut

	if stream := config.GetStreamInfo("AIR_UAT"); stream != nil {
		rows := database.GetStreamProcessedInOut(stream)

		if rows != nil {
			for rows.Next() {
				rows.StructScan(&totalProcessedInOut)
				fmt.Printf("%#v", totalProcessedInOut)
			}
		}
	}

}
