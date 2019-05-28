# EMM Statistics Report Tool
Tool for extracting statistical reports from Ericsson Multi Mediation platform. It uses JSON configuration file to describe EMM resources including the below:
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

In order to use the tool, just copy `emmstats` and `emm-config.yaml` to the same directory, and update the permission of the binary file `emmstats` to be able to run below command:

`./emmstats --help`

This will generate the below output:

```
NAME:
   emmstats - Tool to generate EMM throughput and performance statistic reports.

USAGE:
   emmstats [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     cdrs, c         Input/Output CDRs statistics, cluster name is required
     files, f        Input/Output Files statistics, cluster name is required
     throughput, t   Input/Output Files and CDRs statistics, cluster name is required
     performance, p  CPU and Memory statistics, cluster name is required
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cluster value, --cl value       Name of EMM cluster which contains the logical server
   --lserver. ls value               Name of EMM logical server
   --format value, --fmt value       Output format of the report, valid values (table, csv) (default: "table")
   --start-time value, --sd value    Start time of the report in the format YYMMDDHH24MISS (default: "20190101000000")
   --end-time value, --ed value      End time of the report in the format YYMMDDHH24MISS (default: "20190528162228")
   --ls-database value, --ldb value  Name of adhoc logical server database to specify in CLI without configuring it in EMM config file
   --pf-database value, --pdb value  Name of adhoc performance database to specify in CLI without configuring it in EMM config file
   --db-ip value, --ip value         IP of the adhoc database
   --db-port value, -p value         Port of the adhoc database
   --help, -h                        show help
   --version, -v                     print the version
```

## Sample Commands


## Sample Configuration File

Below is sample `emm-info.json` file which contains description of EMM resources. It must be put on the same directory as the binary file `emm-statistics-report.bin`

```
clusters:
  - name: Test # Cluster Name
    username: mmsuper # Default username used to access logical servers databases
    password: mediation
    logical-servers:
      - name: Server5
        ip: 10.135.5.81
        username: mmsuper
        passwod: thule
        port: 5432
  - name: ryd2
    username: mmsuper
    password: mediation
    logical-servers:
      - name: Server1
        ip: 10.135.3.125
  - name: dev
    username: mmsuper
    password: mediation
    logical-servers:
      - name: Server11
        ip: localhost
        port: 5432
        username: mmsuper
        password: mediation
        database: fm_db_Server11

configurations:
  - name: UAT_Test
    coll-names: ["INPUT", "Output"]
    dist-names: ["BI", "RA"]
    assigned-logical-server:
      name: Server1
      cluster: ryd2
  - name: 4GLTE_INPUT_CDRs
    coll-names: ["to4G_LTE_in_RD", "to4G_LTE_in_RD", "to4G_LTE_in_RD", "to4G_LTE_in_RD"]
    assigned-logical-server:
      name: Server11
      cluster: dev
  - name: HWPGW_INPUT_CDRs
      coll-names: ["toHW_PGW_in_RD", "toHW_PGW_in_JD", "toHW_PGW_in_JE", "toHW_PGW_in_JE"]
      coll-ids: ["14025"]
      assigned-logical-server:
        name: Server11
        cluster: dev
```

