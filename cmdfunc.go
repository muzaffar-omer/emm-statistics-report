package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
	"strings"
	"time"
)

const spinnerUpdateFreq = 1000 * time.Millisecond

// throughput reports the total processed input/output for a logical server, or for specific stream running in a logical
// server
func throughput(context *cli.Context) error {

	s := spinner.New(spinner.CharSets[35], spinnerUpdateFreq) // Build our new spinner

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

			s.Prefix = fmt.Sprintf("%s Stream Throughput ", stream.Name)
			s.Start()

			logicalServer := findLogicalServer(stream.LogicalServer.Name, stream.LogicalServer.Cluster)

			params := AudittrailLogEntryQueryParameters{
				TimeFormat:   chooseGroupByFormat(context.String("group-by")),
				StartTime:    startTimeArg,
				EndTime:      endTimeArg,
				InnodeNames:  stream.CollectorNames,
				InnodeIds:    stream.CollectorIds,
				OutnodeNames: stream.DistributorNames,
				OutnodeIds:   stream.DistributorIds,
			}

			query := parseTemplate("throughput", streamThroughputQueryTemplate, params)

			logger.WithFields(logrus.Fields{
				"command": "throughput",
				"stream":  stream.Name,
				"query":   query,
			}).Debug("Stream throughput query")

			rows := executeQuery(logicalServer, query)

			if rows != nil {
				s.Stop()
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

		s.Prefix = fmt.Sprintf("%s Logical Server Throughput ", logicalServer.Name)
		s.Start()

		params := AudittrailLogEntryQueryParameters{
			TimeFormat: chooseGroupByFormat(context.String("group-by")),
			StartTime:  startTimeArg,
			EndTime:    endTimeArg,
		}

		query := parseTemplate("throughput", lsThroughputQueryTemplate, params)

		logger.WithFields(logrus.Fields{
			"command":        "throughput",
			"logical_server": logicalServer.Name,
			"query":          query,
		}).Debug("Logical server throughput query")

		rows := executeQuery(logicalServer, query)

		if rows != nil {
			s.Stop()
			printResultTable(rows, fmt.Sprintf("Logical Server Throughput : %s", logicalServer.Name))
		}

	} else {
		logger.WithFields(logrus.Fields{
			"command": "throughput",
		}).Fatalln("Invalid command options")
	}

	return nil
}

func cdrs(context *cli.Context) error {
	return nil
}

func validateThroughput(context *cli.Context) error {

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

	return nil
}

func initializeAndValidateGFlags(context *cli.Context) error {
	lsDbname := context.String("ls-dbname")
	pfDbname := context.String("pf-dbname")

	dbIp := context.String("db-ip")
	dbPort := context.String("db-port")

	lserver := context.String("lserver")
	cluster := context.String("cluster")
	stream := context.String("stream")

	// Make sure Adhoc options are not combined with configuration file based options
	if len(lsDbname) > 0 || len(pfDbname) > 0 {
		if len(lserver) > 0 {
			return cli.Exit("Cannot combine --lserver option with adhoc query options", -1)
		}

		if len(cluster) > 0 {
			return cli.Exit("Cannot combine --cluster option with adhoc query options", -1)
		}

		if len(stream) > 0 {
			return cli.Exit("Cannot combine --stream option with adhoc query options", -1)
		}

		if len(dbIp) == 0 {
			return cli.Exit("Missing mandatory option (--db-ip) for adhoc query", -1)
		}

		if len(dbPort) == 0 {
			return cli.Exit("Missing mandatory option (--db-port) for adhoc query", -1)
		}
	}

	// Validate date formats
	startTime := context.String("start-time")
	endTime := context.String("end-time")

	if len(startTime) > 0 {
		_, err := time.Parse(timeFlagFormat, startTime)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Invalid start-time format %s", startTime), -1)
		}
	}

	if len(endTime) > 0 {
		_, err := time.Parse(timeFlagFormat, endTime)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Invalid end-time format %s", endTime), -1)
		}
	}

	// Validate output file format
	outputFormat := context.String("format")
	if len(outputFormat) > 0 && strings.ToLower(outputFormat) != "csv" && strings.ToLower(outputFormat) != "table" {
		return cli.Exit(fmt.Sprintf("Invalid output format %s", outputFormat), -1)
	}

	// Parse EMM configuration file
	emmConfig = parseEMMConfig()

	return nil
}

func validateCdrs(context *cli.Context) error {
	return nil
}

func validateFiles(context *cli.Context) error {
	return nil
}

func validatePerformance(context *cli.Context) error {
	return nil
}

func chooseGroupByFormat(groupByPeriod string) string {
	var groupByFormat string

	switch groupByPeriod {
	case "month":
		groupByFormat = month
		break
	case "day":
		groupByFormat = day
		break
	case "hour":
		groupByFormat = hour
		break
	case "minute":
		groupByFormat = minute
		break
	default:
		groupByFormat = day
	}

	return groupByFormat
}