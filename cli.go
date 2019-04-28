package main

import (
	//"fmt"

	"gopkg.in/urfave/cli.v2"
	"time"
)

//######################### Global Commands ##################################
// Command to generate the CDRs statistics, it could generate the below:
// - Input CDRs for a single server, or all servers
// - Output CDRs for a single server, or all servers
// - Input/Output CDRs for a single server, or all servers
var cdrsCommand = &cli.Command{
	Name:        "cdrs",
	Aliases:     []string{"c"},
	Description: "Input/Output CDRs statistics, cluster name is required",
}

// Command to generate the Files statistics, it could generate the below:
// - Input Files for a single server, or all servers
// - Output Files for a single server, or all servers
// - Input/Output Files for a single server, or all servers
var filesCommand = &cli.Command{
	Name:    "files",
	Aliases: []string{"f"},
	Usage:   "Input/Output Files statistics, cluster name is required",
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

//######################### Adhoc Database Global Flags ##################################
var lsDatabaseGFlag = &cli.StringFlag{
	Name:  "ls-database, ldb",
	Usage: "Name of adhoc logical server database to specify in CLI without configuring it in EMM config file",
}

var perfDatabaseGFlag = &cli.StringFlag{
	Name:  "pf-database, pdb",
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
			outputFormatGFlag,
			startTimeGFlag,
			endTimeGFlag,
			lsDatabaseGFlag,
			perfDatabaseGFlag,
			dbIPGFlag,
			dbPortGFlag,
		},

		Commands: []*cli.Command{
			cdrsCommand,
			filesCommand,
			throughputCommand,
			performanceCommand,
		},
	}

}

// Returns current time formatted in the format YYYYMMDDHH24MISS
func currentTime() string {
	return time.Now().Format("20060102150405")
}

type CliError struct{}

func (e CliError) Error() string {
	return ""
}
