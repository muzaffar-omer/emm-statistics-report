package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
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

var (
	FileConfig EMMFileConfig
	CmdConfig  CmdArgs
)

var logger = logrus.New()

const CONFIG_FILE_NAME = "emm-info.json"

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

					if fileCfg.Clusters[index].LogicalServers[lsIndex].ActiveStream == "" {
						logger.WithFields(logrus.Fields{
							"cluster":        fileCfg.Clusters[index].Name,
							"logical_server": fileCfg.Clusters[index].LogicalServers[lsIndex].Name,
						}).Warn("Missing 'active_stream' field")
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

							fileCfg.Clusters[index].LogicalServers[lsIndex].Username = fileCfg.Clusters[index].DefaultUsername
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

							fileCfg.Clusters[index].LogicalServers[lsIndex].Password = fileCfg.Clusters[index].DefaultPassword
						}
					}
				}
			}

		}
	}

	return true
}

type LogicalServer struct {
	Name         string `json:"name"`
	Ip           string `json:"ip"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         string `json:"port"`
	Database     string `json:"database"`
	ActiveStream string `json:"active_stream"`
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
	ip         string
	username   string
	password   string
	port       string
	reportType string
	days       int8
	dateRange  string
	outputFile string
	outputDir  string
	logLevel   string
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

func (cmdCfg CmdArgs) String() string {
	return fmt.Sprintf("IP:%s\nPort:%s\nUsername:%s\nPassword:%s\n", cmdCfg.ip,
		cmdCfg.port,
		cmdCfg.username,
		cmdCfg.password)
}

func (cmdCfg CmdArgs) LogLevel() string {
	return cmdCfg.logLevel
}

func (cfg *CmdArgs) Parse() {
	flag.StringVar(&cfg.ip, "ip", "localhost", "Postgresql DB instance IP address")
	flag.StringVar(&cfg.username, "username", "mmsuper", "DB user name")
	flag.StringVar(&cfg.password, "password", "thule", "DB user password")
	flag.StringVar(&cfg.port, "port", "5432", "DB port")
	flag.StringVar(&cfg.logLevel, "log_level", "Info", "Sets the logging level, [Debug, Info, Warn, Error, Fatal]")

	flag.Parse()
}

type Stream struct {
	Collectors   []string `json:"collectors"`
	Distributors []string `json:"distributors"`
	Name         string   `json:"name"`
}

type Cluster struct {
	Name            string          `json:"name"`
	DefaultUsername string          `json:"default_username"`
	DefaultPassword string          `json:"default_password"`
	LogicalServers  []LogicalServer `json:"logical_servers"`
}

type EMMFileConfig struct {
	Streams  []Stream  `json:"streams"`
	Clusters []Cluster `json:"clusters"`
}

func FindLsRunningStream(stream *Stream) *LogicalServer {

	for _, cluster := range FileConfig.Clusters {
		for _, ls := range cluster.LogicalServers {
			if ls.ActiveStream == stream.Name {
				return &ls
			}
		}
	}

	return nil
}

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
