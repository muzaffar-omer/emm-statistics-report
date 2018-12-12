package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

/*
	Contains the configuration parameters passed in the command line, also includes the parameters defined in the
	configuration file:
	- connection details:
		- DB IP
		- username
		- password
		- port
	- specify server details JSON files
	- report type
		- input cdrs
		- input files
		- output cdrs
		- output files
	- days
	- daterange
	- output file
	- output directory
	- keep configuration in JSON file:
		- list of logical servers
		- output directory
		- output file
*/

const CONFIG_FILE_NAME = "emm-info.json"
const STREAM_MAP_FORMAT = "(\\w+)@(\\w+):(\\w+)"

var (
	FileConfig EMMFileConfig // contains objects parsed from emm-info.json configuration file
	CmdConfig  CmdArgs       // contains the possible command line arguments that could be provided by the user
)

var logger = logrus.New()

func Log() *logrus.Logger {
	return logger
}

func init() {

	var configuredLevel logrus.Level

	logger.Debug("Just called init() of the configuration package ....")

	// Parse command line args
	CmdConfig.Parse()

	// Convert the level from string into logrus.Level
	configuredLevel, err := logrus.ParseLevel(CmdConfig.logLevel)

	// If there is any error during parsing of the log level, use the default "Info" level
	if err != nil {
		logger.Error("Error parsing the provided cmd log level, will use the default \"Info\" level")
		configuredLevel = logrus.InfoLevel
	}

	logger.SetLevel(configuredLevel)

	logger.WithFields(logrus.Fields{
		"level": configuredLevel,
	}).Debug("Setting log level")

	// Parse EMM configuration file
	jsonFile, err := os.Open(CONFIG_FILE_NAME)

	defer jsonFile.Close()

	// Error while reading the configuration file
	if err != nil {

		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Opening emm-info.json file")

	} else {
		jsonByteArr, err := ioutil.ReadAll(jsonFile)

		logger.WithFields(logrus.Fields{
			"emm-info.json": string(jsonByteArr),
		}).Debug("View contents of emm-info.json")

		// Error during parsing of file contents
		if err != nil {

			logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("Reading contents of emm-info.json file")

		} else {

			// Un packing the configuration file into the Config.fileConfig parameter
			err = json.Unmarshal(jsonByteArr, &FileConfig)

			//fmt.Printf("%#v\n", FileConfig)

			// Validate configuration
			FileConfig.validate()

			logger.WithFields(logrus.Fields{
				"config-struct": fmt.Sprintf("%#v", FileConfig),
			}).Debug("Parsed FileConfig structure")

			if err != nil {
				logger.WithFields(logrus.Fields{
					"error": err,
				}).Error("Parsing contents of emm-info.json file")
			}
		}
	}
}

func isValidStreamMapFormat(streamMap string) bool {
	// Stream Map Format : <Stream Name>@<Cluster Name>:<Logical Server Name>
	return regexp.MustCompile(STREAM_MAP_FORMAT).MatchString(streamMap)
}

// Parses a stream mapping and extracts the streamName, clusterName, logicalServerName
// from the parsed string
func extractStreamMapParams(streamMap string) (streamName, clusterName, logicalServerName string) {
	destructuredStreamMap := regexp.MustCompile(STREAM_MAP_FORMAT).FindAllStringSubmatch(streamMap, -1)

	if isValidStreamMapFormat(streamMap) && len(destructuredStreamMap) == 1 && len(destructuredStreamMap[0]) > 1 {
		streamName = destructuredStreamMap[0][1]        // Stream Name
		clusterName = destructuredStreamMap[0][2]       // Cluster Name
		logicalServerName = destructuredStreamMap[0][3] // Logical Server Name
	} else {
		streamName, clusterName, logicalServerName = "", "", ""
	}

	return
}

// Validates the content of the emm-info.json file, each part of the configuration file has
// specific validation rules
func (fileCfg *EMMFileConfig) validate() bool {

	// Validate streams
	// 1 - Check there are streams defined
	if len(fileCfg.Streams) == 0 {
		logger.Error("No streams defined in " + CONFIG_FILE_NAME + ", check 'name' field is defined under" +
			" main configuration file structure")
	}

	// Validate the fields of each stream
	for _, stream := range fileCfg.Streams {

		// Check missing stream names
		if stream.Name == "" {
			logger.Error("Missing stream name, check that 'name' field is defined for all streams")
		}
	}

	// Validate stream mapping
	for index, streamMap := range fileCfg.StreamMapping {

		if !isValidStreamMapFormat(streamMap) {
			logger.WithFields(logrus.Fields{
				"stream_map": streamMap,
			}).Error("Invalid stream mapping, will be ignored")

			// Remove streamMap from StreamMapping array
			fileCfg.StreamMapping = append(fileCfg.StreamMapping[:index], fileCfg.StreamMapping[index+1:]...)
		}
	}

	// Validate clusters information
	if len(fileCfg.Clusters) == 0 {
		logger.Error("No clusters defined in " + CONFIG_FILE_NAME + ", check that 'clusters' field is defined" +
			"under the main configuration file structure")
	} else {
		for index, _ := range fileCfg.Clusters {

			// Validate cluster name, and default username and password
			if fileCfg.Clusters[index].Name == "" {
				logger.Error("Missing cluster 'name' field, check that all clusters have 'name' field")
			}

			if fileCfg.Clusters[index].DefaultUsername == "" {
				logger.WithFields(logrus.Fields{
					"cluster": fileCfg.Clusters[index].Name,
				}).Warn("Missing 'default_username' field in cluster definition")
			}

			if fileCfg.Clusters[index].DefaultPassword == "" {
				logger.WithFields(logrus.Fields{
					"cluster": fileCfg.Clusters[index].Name,
				}).Warn("Missing 'default_password' field in cluster definition")
			}

			// Check definitions of logical servers
			if len(fileCfg.Clusters[index].LogicalServers) == 0 {
				logger.WithFields(logrus.Fields{
					"cluster": fileCfg.Clusters[index].Name,
				}).Warn("Missing 'logical_servers' field in cluster definition")
			} else {

				for lsIndex, _ := range fileCfg.Clusters[index].LogicalServers {

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Name == "" {
						logger.WithFields(logrus.Fields{
							"cluster": fileCfg.Clusters[index].Name,
						}).Error("Missing 'name' field, check that all logical servers have " +
							"this field")
					}

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Ip == "" {
						logger.WithFields(logrus.Fields{
							"cluster":        fileCfg.Clusters[index].Name,
							"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
						}).Error("Missing 'ip' field")
					}

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Database == "" {
						logger.WithFields(logrus.Fields{
							"cluster":        fileCfg.Clusters[index].Name,
							"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
						}).Error("Missing 'database' field")
					}

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Port == "" {
						logger.WithFields(logrus.Fields{
							"cluster":        fileCfg.Clusters[index].Name,
							"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
						}).Error("Missing 'port' field")
					}

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Username == "" {
						if fileCfg.Clusters[index].DefaultUsername == "" {
							logger.WithFields(logrus.Fields{
								"cluster":        fileCfg.Clusters[index].Name,
								"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
							}).Error("Missing 'username' field, and missing cluster 'default_username'")
						} else {
							logger.WithFields(logrus.Fields{
								"cluster":        fileCfg.Clusters[index].Name,
								"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
							}).Warn("Missing 'username' field, will use cluster 'default_username'")

							fileCfg.Clusters[index].LogicalServers[lsIndex].Username =
								fileCfg.Clusters[index].DefaultUsername
						}
					}

					if fileCfg.Clusters[index].LogicalServers[lsIndex].Password == "" {
						if fileCfg.Clusters[index].DefaultPassword == "" {
							logger.WithFields(logrus.Fields{
								"cluster":        fileCfg.Clusters[index].Name,
								"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
							}).Error("Missing 'password' field, and missing cluster 'default_password'")
						} else {
							logger.WithFields(logrus.Fields{
								"cluster":        fileCfg.Clusters[index].Name,
								"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
							}).Warn("Missing 'password' field, will use cluster 'default_password'")

							fileCfg.Clusters[index].LogicalServers[lsIndex].Password =
								fileCfg.Clusters[index].DefaultPassword
						}
					}
				}
			}

		}
	}

	return true
}

// looks in the list of defined logical servers using logical server name, and cluster name, and returns
// logical server object
func (fileCfg EMMFileConfig) getLogicalServer(clusterName string, logicalServerName string) *LogicalServer {

	logger.WithFields(logrus.Fields{
		"logical_name": logicalServerName,
		"cluster_name": clusterName,
	}).Debug("Looking for logical server")

	for _, cluster := range FileConfig.Clusters {

		if cluster.Name == clusterName {

			logger.WithFields(logrus.Fields{
				"logical_name": logicalServerName,
				"cluster_name": clusterName,
			}).Debug("Found the cluster")

			for _, ls := range cluster.LogicalServers {
				if ls.Name == logicalServerName {

					logger.WithFields(logrus.Fields{
						"logical_name": logicalServerName,
						"cluster_name": clusterName,
					}).Debug("Found the logical server")

					return &ls
				}
			}
		}

	}

	return nil
}

type LogicalServer struct {
	Name     string `json:"name"`
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

func (this LogicalServer) Equals(another *LogicalServer) bool {
	if this.Username == another.Username &&
		this.Ip == another.Ip &&
		this.Password == another.Password &&
		this.Port == another.Port {
		return true
	}

	return false
}

type CmdArgs struct {
	ip            string
	username      string
	password      string
	port          string
	reportType    string
	days          int8
	dateRange     string
	outputFile    string
	outputDir     string
	logLevel      string
	groupBy       string
	fromDate      string
	toDate        string
	numberOfDay   string
	outputFormat  string
	stream        string
	operationType int
}

func (cmdCfg CmdArgs) Ip() string {
	return cmdCfg.ip
}

func (cmdCfg CmdArgs) Username() string {
	return cmdCfg.username
}

func (cmdCfg CmdArgs) Password() string {
	return cmdCfg.password
}

func (cmdCfg CmdArgs) Port() string {
	return cmdCfg.port
}

func (cmdCfg CmdArgs) Stream() string {
	return cmdCfg.stream
}

func (cmdCfg CmdArgs) GroupBy() string {
	return cmdCfg.groupBy
}

func (cmdCfg CmdArgs) String() string {
	return fmt.Sprintf("IP:%s\nPort:%s\nUsername:%s\nPassword:%s\n", cmdCfg.ip,
		cmdCfg.port,
		cmdCfg.username,
		cmdCfg.password)
}

func (cmdCfg CmdArgs) LogLevel() string {
	return cmdCfg.logLevel
}

func (cmdCfg CmdArgs) OperationType() int {
	return cmdCfg.operationType
}

func (cmdCfg CmdArgs) FromDate() string {
	return cmdCfg.fromDate
}

func (cmdCfg CmdArgs) ToDate() string {
	return cmdCfg.toDate
}

func (cfg *CmdArgs) Parse() {
	lastDay := time.Unix(time.Now().Unix()-(24*60*60), 0)

	flag.StringVar(&cfg.ip, "ip", "localhost", "Postgresql DB instance IP address")
	flag.StringVar(&cfg.username, "username", "mmsuper", "DB user name")
	flag.StringVar(&cfg.password, "password", "thule", "DB user password")
	flag.StringVar(&cfg.port, "port", "5432", "DB port")
	flag.StringVar(&cfg.logLevel, "log-level", "Error", "Sets the logging level, [Debug, Info, "+
		"Warn, Error, Fatal]")
	flag.StringVar(&cfg.groupBy, "group-by", "day", "Specifies the intervals for grouping of the "+
		"result [minute, hour, day, month], default value is 'day'")
	flag.StringVar(&cfg.fromDate, "from-date", "19700101", "Specifies the start date for generation "+
		"of the report in the format YYYYMMDD")
	flag.StringVar(&cfg.toDate, "to-date", lastDay.Format("20060102"), "Specifies the end date for generation "+
		"of the report in the format YYYYMMDD")
	flag.StringVar(&cfg.outputFormat, "output-format", "table", "Specifies the format of the result [table, csv]")
	flag.StringVar(&cfg.stream, "stream", "", "Stream name defined in the EMM configuration file")
	flag.IntVar(&cfg.operationType, "query-type", 1, "Specifies the required type of query "+
		"(operation), below are the possible values:\n"+
		"1 - Stream processed input/output grouped by minute, hour, day, or month, it requires the group-by "+
		"parameter to be specified (default group-by value is day)")

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(-1)
	}

	// Validate command line args
}

type Stream struct {
	Collectors   []string `json:"collectors"`
	Distributors []string `json:"distributors"`
	Name         string   `json:"name"`
}

type Cluster struct {
	Name string `json:"name"`
	// will be used in case no username defined for the logical server
	DefaultUsername string `json:"default_username"`

	// will be used in case no password defined for the logical server
	DefaultPassword string `json:"default_password"`

	// list of all logical servers details including IP, db name, logical server name, username, password
	LogicalServers []LogicalServer `json:"logical_servers"`
}

type EMMFileConfig struct {
	Streams       []Stream  `json:"streams"`        // definition of collectors/distributors for each stream
	Clusters      []Cluster `json:"clusters"`       // definition of all clusters and their logical servers details
	StreamMapping []string  `json:"stream_mapping"` // list of streams mapped to logical servers and their clusters
}

// Looks in the stream_mapping defined in the configuration file
// and finds the logical server which is running Stream based on
// streamName
func FindLsRunningStream(stream *Stream) *LogicalServer {

	logger.WithFields(logrus.Fields{
		"stream": stream.Name,
	}).Debug("Looking for the logical server where this stream is assigned")

	for _, streamMap := range FileConfig.StreamMapping {
		streamName, clusterName, logicalServerName := extractStreamMapParams(streamMap)

		if streamName == stream.Name {

			logger.WithFields(logrus.Fields{
				"stream":         stream.Name,
				"stream_mapping": streamMap,
			}).Debug("Found stream mapping definition")

			return FileConfig.getLogicalServer(clusterName, logicalServerName)
		}
	}

	return nil
}

// Looks in the streams defined in the configuration file, and returns the
// Stream object matching the streamName
func GetStreamInfo(streamName string) *Stream {

	logger.WithFields(logrus.Fields{
		"streams": len(FileConfig.Streams),
	}).Debug("Number of configured streams")

	for _, stream := range FileConfig.Streams {
		if stream.Name == streamName {
			return &stream
		}
	}

	logger.WithFields(logrus.Fields{
		"stream": streamName,
	}).Error("Stream not found")

	return nil
}
