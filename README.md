# EMM Statistics Report Tool
Tool for extracting statistical reports from Ericsson Multi Mediation platform. It depends on the JSON file which contains the definition of EMM resources including the below:
1. Logical Servers:
    * Name
    * Virtual IP (in case of clustered environment), in case of standalone deployments, it will be the same IP as the EMM node IP 
    * Database port
    * Database instance name
2. Clusters:
    * Name
    * Assigned logical servers
    * Default username and password for database access
3. Streams (representas EMM configurations):
    * Name, it is independent from the actual EMM configuration name used in the platform. It used to specify only to identify the configuration within the JSON file.
    * Collectors configured in the configuration
    * Distributors configured in the configuration
4. Streams Mapping, contains the mapping between stream names and their assigned logical servers

## Installation

In order to use the tool, just copy `emm-statistics-report.bin` and `emm-info.json` to the same directory, and update the permission of the binary file `emm-statistics-report.bin` to be able to run the command using the below:

`./emm-statistics-report.bin --help`

This will generate the below output:

```
Usage of ./emm-statistics-report:
  -from-date string
    	Specifies the start date for generation of the report in the format YYYYMMDD (default "19700101")
  -group-by string
    	Specifies the intervals for grouping of the result [minute, hour, day, month], default value is 'day' (default "day")
  -ip string
    	Postgresql DB instance IP address (default "localhost")
  -log-level string
    	Sets the logging level, [Debug, Info, Warn, Error, Fatal] (default "Error")
  -ls string
    	Logical server name in format Server1@RYD1
  -output-format string
    	Specifies the format of the result [table, csv] (default "table")
  -password string
    	DB user password (default "thule")
  -port string
    	DB port (default "5432")
  -query-type int
    	Specifies the required type of query (operation), below are the possible values:
    	1 - Stream processed input/output grouped by minute, hour, day, or month, it requires the group-by parameter to be specified (default group-by value is day)2 - Logical server processed input/output grouped by minute, hour, day, or month, it requires the group-by parameter to be specified (default group-by value is day), requires setting --ls parameter (default 1)
  -stream string
    	Stream name defined in the EMM configuration file
  -to-date string
    	Specifies the end date for generation of the report in the format YYYYMMDD (default "20190205")
  -username string
    	DB user name (default "mmsuper")
```

## Use Cases

1. Generate statistics for daily

`./emm-statistics-report.bin --stream "MSS_ALL_INPUT"`

2. Generate daily statistics starting from 1st Jan 2019 till yesterday

`./emm-statistics-report.bin --stream "MSS_ALL_INPUT" --from-date="20190101"`

2. Generate daily statistics starting from 1st Jan 2019 till 1st Feb 2019

`./emm-statistics-report.bin --stream "MSS_ALL_INPUT" --from-date="20190101" --to-date="20190201"`

3. Generate hourly statistics starting from 1st Jan 2019 till yesterady

`./emm-statistics-report.bin --stream "MSS_ALL_INPUT" --from-date="20190101" --group-by="hour"`

3. Generate statistics aggregated by month starting from 1st Jan 2019 till yesterady

`./emm-statistics-report.bin --stream "MSS_ALL_INPUT" --from-date="20190101" --group-by="month"`

## Sample Configuration File

```
{
	"stream_mapping" : [
		"MSS_ALL_INPUT@RYD1:Server14",
	],
  "streams": [
	{
      "name": "MSS_ALL_INPUT",
      "collectors": [
        "MSS_MAIN"
      ],
      "distributors": [
        "tmp_Z"
      ]
    }
  ],
  "clusters": [
    {
      "default_username": "xxxxxx",
      "default_password": "xxxxx",
      "name": "RYD1",
      "logical_servers": [
	     {
          "name": "Server14",
          "ip": "10.135.3.191",
          "database": "fm_db_Server14",
          "port": "5690"
        }
      ]
    }
  ]
}

```

