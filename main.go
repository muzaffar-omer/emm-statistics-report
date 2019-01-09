package main

import (
	config "emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"emm-statistics-report/stats"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"time"
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
	case 2:
		OperationLogicalServerGroupedProcessedInOut()
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

		var fromDate time.Time
		var toDate time.Time

		logger.WithFields(logrus.Fields{
			"stream_name": config.CmdConfig.Stream(),
		}).Debug("Stream is defined in configuration file")

		tmpDate, err := database.ConvertCmdDateToTime(config.CmdConfig.FromDate())

		if err != nil {
			logger.WithFields(logrus.Fields{
				"from-date": config.CmdConfig.FromDate(),
				"error":     err,
			}).Panic("Could no convert provided date into internal time format")
		} else {
			fromDate = tmpDate
		}

		tmpDate, err = database.ConvertCmdDateToTime(config.CmdConfig.ToDate())

		if err != nil {
			logger.WithFields(logrus.Fields{
				"to-date": config.CmdConfig.FromDate(),
				"error":   err,
			}).Panic("Could no convert provided date into internal time format")
		} else {
			toDate = tmpDate
		}

		rows := database.GetStreamProcessedInOut(stream, config.CmdConfig.GroupBy(), fromDate, toDate)

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

func OperationLogicalServerGroupedProcessedInOut() {
	var totalGroupedProcessedInOut database.TotalGroupedProcessedInOut
	var statisticalRecords []stats.Statistical

	table := tablewriter.NewWriter(os.Stdout)

	logger.WithFields(logrus.Fields{
		"logical_server": config.CmdConfig.LogicalServer(),
	}).Debug("Generating OperationGroupedLogicalServerProcessedInOut report")

	logicalServer := config.GetLogicalServerInfo(config.CmdConfig.LogicalServer())

	if logicalServer != nil {

		var fromDate time.Time
		var toDate time.Time

		tmpDate, err := database.ConvertCmdDateToTime(config.CmdConfig.FromDate())

		if err != nil {
			logger.WithFields(logrus.Fields{
				"from-date": config.CmdConfig.FromDate(),
				"error":     err,
			}).Panic("Could no convert provided date into internal time format")
		} else {
			fromDate = tmpDate
		}

		tmpDate, err = database.ConvertCmdDateToTime(config.CmdConfig.ToDate())

		if err != nil {
			logger.WithFields(logrus.Fields{
				"to-date": config.CmdConfig.FromDate(),
				"error":   err,
			}).Panic("Could no convert provided date into internal time format")
		} else {
			toDate = tmpDate
		}

		rows := database.GetLogicalServerProcessedInOut(logicalServer, config.CmdConfig.GroupBy(), fromDate, toDate)

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
