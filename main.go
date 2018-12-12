package main

import (
	config "emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"emm-statistics-report/stats"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = config.Log()

func main() {

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true})

	switch config.CmdConfig.OperationType() {
	case 1:
		OperationGroupedProcessedInOut()
		break
	default:
		OperationGroupedProcessedInOut()
	}
}

// Possible operations:
// 1 - Get processed input/output grouped by minute, hour, day, or month
func OperationGroupedProcessedInOut() {
	var totalGroupedProcessedInOut database.TotalGroupedProcessedInOut
	var statisticalRecords []stats.Statistical

	table := tablewriter.NewWriter(os.Stdout)

	logger.WithFields(logrus.Fields{
		"stream_name": config.CmdConfig.Stream(),
	}).Debug("Generating OperationGroupedProcessedInOut report")

	if stream := config.GetStreamInfo(config.CmdConfig.Stream()); stream != nil {

		logger.WithFields(logrus.Fields{
			"stream_name": config.CmdConfig.Stream(),
		}).Debug("Stream is defined in configuration file")

		rows := database.GetGroupedStreamProcessedInOut(stream, config.CmdConfig.GroupBy())

		if rows != nil {

			for rows.Next() {
				totalGroupedProcessedInOut = database.TotalGroupedProcessedInOut{}
				rows.StructScan(&totalGroupedProcessedInOut)
				statisticalRecords = append(statisticalRecords, totalGroupedProcessedInOut)
				table.Append(totalGroupedProcessedInOut.AsArray())
			}

			table.SetHeader(totalGroupedProcessedInOut.Header())
			table.Render()

			table = stats.CreateStatisticsTable(statisticalRecords)
			table.Render()
		}
	}
}
