package main

import (
	//"fmt"

	"gopkg.in/urfave/cli.v2"
	"time"
)

const (
	timeFlagFormat = "20060102150405"
	errorExitCode  = -1
)

//######################### Global Commands ##################################
// Command to generate the CDRs statistics, it could generate the below:
// - Input CDRs for a single server, or all servers
// - Output CDRs for a single server, or all servers
// - Input/Output CDRs for a single server, or all servers
var cdrsCommand = &cli.Command{
	Name:        "cdrs",
	Aliases:     []string{"c"},
	Usage: "Input/Output CDRs statistics, cluster name is required",
	Action:      cdrs,
	Before: validateCdrsOptions,
}

// Command to generate the Files statistics, it could generate the below:
// - Input Files for a single server, or all servers
// - Output Files for a single server, or all servers
// - Input/Output Files for a single server, or all servers
var filesCommand = &cli.Command{
	Name:    "files",
	Aliases: []string{"f"},
	Usage:   "Input/Output Files statistics, cluster name is required",
	Before: validateFilesOptions,
}

// Command to generate the Throughput (Files and CDRs) statistics, it could
// generate the below:
// - Input Throughput for a single server, or all servers
// - Output Throughput for a single server, or all servers
// - Input/Output Throughput for a single server or all servers
var throughputCommand = &cli.Command{
	Name:    "throughput",
	Aliases: []string{"t"},
	Usage:   "Input/Output Files and CDRs statistics, cluster name is required",
	Action:  throughput,
	Before: validateThroughputOptions,
}

// Command to generate CPU and Memory statistics as below:
// - For a single server, or all servers
var performanceCommand = &cli.Command{
	Name:    "performance",
	Aliases: []string{"p"},
	Usage:   "CPU and Memory statistics, cluster name is required",
	Subcommands: []*cli.Command{
		cpuCommand,
		memCommand,
	},
	Before: validatePerformanceOptions,
}

//######################### Performance Subcommands ##################################
// Command to generate CPU statistics as below:
// - For a single server, or all servers
var cpuCommand = &cli.Command{
	Name:  "cpu",
	Usage: "CPU statistics, cluster name is required",
}

// Command to generate Memory statistics as below:
// - For a single server, or all servers
var memCommand = &cli.Command{
	Name:    "memory",
	Aliases: []string{"mem"},
	Usage:   "Memory statistics, cluster name is required",
}

//######################### Global Flags ##################################
var clusterGFlag = &cli.StringFlag{
	Name:  "cluster, cl",
	Usage: "Name of EMM cluster which contains the logical server",
}

var logicalServerGFlag = &cli.StringFlag{
	Name:  "lserver, ls",
	Usage: "Name of EMM logical server",
}

var outputFormatGFlag = &cli.StringFlag{
	Name:  "format, fmt",
	Usage: "Output format of the report, valid values (table, csv)",
	Value: "table",
}

var startTimeGFlag = &cli.StringFlag{
	Name:  "start-time, sd",
	Usage: "Start time of the report in the format YYMMDDHH24MISS",
	Value: "20190101000000",
}

var endTimeGFlag = &cli.StringFlag{
	Name:  "end-time, ed",
	Usage: "End time of the report in the format YYMMDDHH24MISS",
	Value: currentTime(),
}

var groupByGFlag = &cli.StringFlag{
	Name:  "group-by, gb",
	Usage: "Time interval for grouping of the result, possible values are minute, hour, day, month",
	Value: "day",
}

var streamGFlag = &cli.StringFlag{
	Name:  "stream, s",
	Usage: "Name of the stream defined in YAML configuration file",
}

var verboseGFlag = &cli.BoolFlag{
	Name:  "verbose, v",
	Usage: "Verbose mode (set log level to debug)",
	Value: false,
}

//######################### Adhoc Database Global Flags ##################################
var lsDatabaseGFlag = &cli.StringFlag{
	Name:  "ls-dbname, ldb",
	Usage: "Name of adhoc logical server database to specify in CLI without configuring it in EMM config file",
}

var perfDatabaseGFlag = &cli.StringFlag{
	Name:  "pf-dbname, pdb",
	Usage: "Name of adhoc performance database to specify in CLI without configuring it in EMM config file",
}

var dbIPGFlag = &cli.StringFlag{
	Name:  "db-ip, ip",
	Usage: "IP of the adhoc database",
}

var dbPortGFlag = &cli.StringFlag{
	Name:  "db-port, p",
	Usage: "Port of the adhoc database",
}

func CreateCliApp() *cli.App {
	return &cli.App{
		Name:  "emmstats",
		Usage: "Tool to generate EMM throughput and performance statistic reports",
		Authors: []*cli.Author{
			{Name: "Muzaffar", Email: "muzaffar.omer@gmail.com"},
		},

		Flags: []cli.Flag{
			clusterGFlag,
			logicalServerGFlag,
			streamGFlag,
			verboseGFlag,
			outputFormatGFlag,
			startTimeGFlag,
			endTimeGFlag,
			lsDatabaseGFlag,
			perfDatabaseGFlag,
			dbIPGFlag,
			dbPortGFlag,
			groupByGFlag,
		},

		Commands: []*cli.Command{
			cdrsCommand,
			filesCommand,
			throughputCommand,
			performanceCommand,
		},
		Before: initializeAndValidateGFlags,
	}

}

// Returns current time formatted in the format YYYYMMDDHH24MISS
func currentTime() string {
	return time.Now().Format("20060102150405")
}
