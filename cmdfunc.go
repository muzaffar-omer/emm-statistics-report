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

	verbose := context.Bool("verbose")

	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

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

			logicalServer := findLogicalServer(stream.LogicalServer.Name, stream.LogicalServer.Cluster)

			groupByDateFormat := context.String("group-by")

			switch groupByDateFormat {
			case "month":
				groupByDateFormat = month
				break
			case "day":
				groupByDateFormat = day
				break
			case "hour":
				groupByDateFormat = hour
				break
			case "minute":
				groupByDateFormat = minute
				break
			default:
				groupByDateFormat = day
			}

			params := AudittrailLogEntryQueryParameters{
				TimeFormat:   groupByDateFormat,
				StartTime:    startTimeArg,
				EndTime:      endTimeArg,
				InnodeNames:  stream.CollectorNames,
				InnodeIds:    stream.CollectorIds,
				OutnodeNames: stream.DistributorNames,
				OutnodeIds:   stream.DistributorIds,
			}

			query := parseTemplate("throughput", streamThroughputQueryTemplate, params)

			logger.WithFields(logrus.Fields{
				"command" : "throughput",
				"stream" : stream.Name,
				"query" : query,
			}).Debug("Stream throughput query")

			rows := executeQuery(logicalServer, query)

			if rows != nil {
				printResultTable(rows, fmt.Sprintf("Stream Throughput : %s", stream.Name))
			}

		} else {
			logger.WithFields(logrus.Fields{
				"command": "throughput",
			}).Fatalf("%s stream is not assigned to any logical server", stream.Name)
		}
	} else if len(logicalServerArg) > 0 && len(clusterArg) > 0 {
		// Generate throughput report for a complete logical server audittraillogentry
		logicalServer := findLogicalServer(logicalServerArg, clusterArg)

		groupByDateFormat := context.String("group-by")

		switch groupByDateFormat {
		case "month":
			groupByDateFormat = month
			break
		case "day":
			groupByDateFormat = day
			break
		case "hour":
			groupByDateFormat = hour
			break
		case "minute":
			groupByDateFormat = minute
			break
		default:
			groupByDateFormat = day
		}

		params := AudittrailLogEntryQueryParameters{
			TimeFormat: groupByDateFormat,
			StartTime:  startTimeArg,
			EndTime:    endTimeArg,
		}

		query := parseTemplate("throughput", lsThroughputQueryTemplate, params)

		logger.WithFields(logrus.Fields{
			"command" : "throughput",
			"logical_server" : logicalServer.Name,
			"query" : query,
		}).Debug("Logical server throughput query")

		rows := executeQuery(logicalServer, query)

		if rows != nil {
			printResultTable(rows, fmt.Sprintf("Logical Server Throughput : %s", logicalServer.Name))
		}

	} else {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Invalid command options")
	}

	return nil
}
