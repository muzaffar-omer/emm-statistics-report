package main

import (
	"database/sql"
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

		rows := database.GetAllRows(session, "fraud_cases")

		var msisdn, file_name, insert_date, imsi, date sql.NullString

		for rows.Next() {

			if err := rows.Scan(&msisdn, &file_name, &insert_date, &imsi, &date); err != nil {
				fmt.Println("Error getting rows")
				fmt.Println(err)
			}

			fmt.Printf("msisdn : %s, file_name : %s, imsi : %s, date : %s \n", msisdn.String, file_name.String, imsi.String, date.String)
		}

		//if rows != nil {
		//	colmuns, _ := rows.;
		//	for column, _ := range colmuns {
		//		fmt.Print(column)
		//	}
		//}
	}
}
