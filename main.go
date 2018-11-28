package main

import (
	config "emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
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

	var totalGroupedProcessedInOut database.TotalGroupedProcessedInOut

	table := tablewriter.NewWriter(os.Stdout)

	if stream := config.GetStreamInfo(config.CmdConfig.Stream()); stream != nil {

		rows := database.GetGroupedStreamProcessedInOut(stream, config.CmdConfig.GroupBy())

		if rows != nil {
			for rows.Next() {
				totalGroupedProcessedInOut = database.TotalGroupedProcessedInOut{}
				rows.StructScan(&totalGroupedProcessedInOut)
				table.Append(totalGroupedProcessedInOut.AsArray())
			}

			table.SetHeader(totalGroupedProcessedInOut.Header())
			table.Render()
		}
	}

}
