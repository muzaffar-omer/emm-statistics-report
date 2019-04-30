package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// throughput reports the total processed input/output for a logical server, or for specific stream running in a logical
// server
func throughput(context *cli.Context) error {

	startTimeArg := context.String("start-time")
	endTimeArg := context.String("end-time")

	// Logical server name, and cluster are required to generate throughput for specific logical server
	logicalServerArg := context.String("lserver")
	clusterArg := context.String("cluster")

	// Stream name is required to generate throughput for specific stream
	streamArg := context.String("stream")

	// Stream and logical server information are exclusive, it is not possible to specify both, either specify
	// logical server details (i.e. logical server, and cluster name). Or specify stream name only
	if (len(logicalServerArg) > 0 || len(clusterArg) > 0) && len(streamArg) > 0 {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Stream and Logical Server flags cannot be specified at the same time in CLI global options" +
			". Either specify a stream, or specify logical server information (i.e. logical server name, and cluster " +
			"name")
	} else if len(logicalServerArg) > 0 && len(clusterArg) == 0 {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Cluster name is missing")
	} else if len(logicalServerArg) == 0 && len(clusterArg) > 0 {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Logical server name is missing")
	}

	// Generate throughput report for a stream
	if len(streamArg) > 0 {

		stream := findStream(streamArg)

		if stream.LogicalServer != nil {
			// logicalServer := findLogicalServer(stream.LogicalServer.Name, stream.LogicalServer.Cluster)

		} else {
			logger.WithFields(logrus.Fields{
				"command": "throughput",
			}).Fatalf("%s stream is not assigned to any logical server", stream.Name)
		}
	} else if len(logicalServerArg) > 0 && len(clusterArg) > 0 {
		logicalServer := findLogicalServer(logicalServerArg, clusterArg)
		params := AudittrailLogEntryQueryParameters{
			TimeFormat: defaultDBTimeFormat,
			StartTime:  startTimeArg,
			EndTime:    endTimeArg,
		}
		query := parseTemplate("throughput", throughputQueryTemplate, params)

		rows := executeQuery(logicalServer, query)

		if rows != nil {
			var totalGroupedProcessedInOut TotalGroupedProcessedInOut

			for rows.Next() {
				totalGroupedProcessedInOut = TotalGroupedProcessedInOut{}
				rows.StructScan(&totalGroupedProcessedInOut)

				fmt.Printf("%s,%d,%d,%d\n", totalGroupedProcessedInOut.Time,
					totalGroupedProcessedInOut.TotalInputBytes,
					totalGroupedProcessedInOut.TotalInputCdrs,
					totalGroupedProcessedInOut.TotalInputFiles)
			}

		} else {
			fmt.Println("Rows is nil")
		}

	} else {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Invalid command options")
	}

	return nil
}
