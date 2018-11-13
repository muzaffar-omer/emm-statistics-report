package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
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

type Config struct {
	fileConfig EMMFileConfig
	cmdConfig  CmdArgs
}

func (cfg *Config) Init() {

	// Parse command line args
	cfg.cmdConfig.Parse()

	// Parse EMM configuration file
	jsonFile, err := os.Open("emm-info.json")

	// Error while reading the configuration file
	if err != nil {

		log.WithFields(log.Fields{
			"error": err,
		}).Error("Opening emm-info.json file")

	} else {
		jsonByteArr, err := ioutil.ReadAll(jsonFile)

		// Error during parsing of file contents
		if err != nil {

			log.WithFields(log.Fields{
				"error": err,
			}).Error("Reading contents of emm-info.json file")

		} else {

			// Un packing the configuration file into the Config.fileConfig parameter
			err = json.Unmarshal([]byte(jsonByteArr), &cfg.fileConfig)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Parsing contents of emm-info.json file")
			}
		}
	}
}

func (cfg Config) Ip() string {
	return cfg.cmdConfig.ip
}

func (cfg Config) Username() string {
	return cfg.cmdConfig.username
}

func (cfg Config) Password() string {
	return cfg.cmdConfig.password
}

func (cfg Config) Port() string {
	return cfg.cmdConfig.port
}

func (cfg Config) Clusters() []Cluster {
	return cfg.fileConfig.Clusters
}

func (cfg Config) Streams() []Stream {
	return cfg.fileConfig.Streams
}

func (cfg Config) String() string {
	return fmt.Sprintf("IP:%s\nPort:%s\nUsername:%s\nPassword:%s\n", cfg.cmdConfig.ip,
		cfg.cmdConfig.port,
		cfg.cmdConfig.username,
		cfg.cmdConfig.password)
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
}

func (cfg *CmdArgs) Parse() {
	flag.StringVar(&cfg.ip, "ip", "localhost", "Postgresql DB instance IP address")
	flag.StringVar(&cfg.username, "username", "mmsuper", "DB user name")
	flag.StringVar(&cfg.password, "password", "thule", "DB user password")
	flag.StringVar(&cfg.port, "port", "5432", "DB port")

	flag.Parse()
}
