package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// defaultEMMConfigFile contains the default name of the EMM YAML configuration file
const defaultEMMConfigFile = "emm-config.yaml"

// emmConfig contains the parsed EMM YAML configuration file
var emmConfig *Config

// Config represents all the modules and submodules of the EMM YAML configuration file
type Config struct {
	Clusters []*Cluster `yaml:"clusters"`
	Streams  []*Stream  `yaml:"configurations"`
}

//######################### Main Modules ##################################

// Stream represents EMM business logic, it specifies the names of collectors and distributors to use in queries and
// specifies the logical server where the stream is running. Name of stream is independent from the name of the business
// logic used in production EMM. It is just a name
type Stream struct {
	Name             string                 `yaml:"name"`
	CollectorNames   []string               `yaml:"coll-names"`
	DistributorNames []string               `yaml:"dist-names"`
	CollectorIds     []string               `yaml:"coll-ids"`
	DistributorIds   []string               `yaml:"dist-ids"`
	LogicalServer    *AssignedLogicalServer `yaml:"assigned-logical-server"`
}

// Cluster is the top-level modules which contains the definition of the logical servers
type Cluster struct {
	Name           string           `yaml:"name"`
	Username       string           `yaml:"username"`
	Password       string           `yaml:"password"`
	LogicalServers []*LogicalServer `yaml:"logical-servers"`
}

//######################### Sub-Modules ##################################

// AssignedLogicalServer is a sub-module used in definition of streams, it specifies the name of the logical server, and
// the name of the cluster where the stream is running
type AssignedLogicalServer struct {
	Name    string `yaml:"name"`
	Cluster string `yaml:"cluster"`
}

// LogicalServer is a sub-module used in the Cluster top-level module, it specifies all the properties of the logical
// server
type LogicalServer struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

// Equals compares the current logical server with another logical server
func (l LogicalServer) Equals(another *LogicalServer) bool {
	if l.Username == another.Username &&
		l.IP == another.IP &&
		l.Password == another.Password &&
		l.Port == another.Port {
		return true
	}

	return false
}

// findStream Looks in the streams defined in the configuration file, and returns the Stream object matching the
// streamName
func (c Config) FindStream(streamName string) *Stream {

	logger.WithFields(logrus.Fields{
		"stream": streamName,
	}).Debug("Searching for stream in EMM configuration")

	for _, stream := range c.Streams {
		if stream.Name == streamName {
			return stream
		}
	}

	return nil
}

// findLogicalServer searches for a logical server in EMM configuration using logical server name, and cluster name
func (c Config) FindLogicalServer(logicalServerName string, clusterName string) *LogicalServer {

	logger.WithFields(logrus.Fields{
		"logical-server-name": logicalServerName,
		"cluster-name":        clusterName,
	}).Debug("Searching for logical server in EMM configuration")

	cluster := c.FindCluster(clusterName)

	if cluster != nil {
		for _, logicalServer := range cluster.LogicalServers {
			if logicalServer.Name == logicalServerName {
				return logicalServer
			}
		}
	}

	return nil
}

// findCluster searches for a cluster definition in EMM configuration file using cluster name
func (c Config) FindCluster(clusterName string) *Cluster {

	logger.WithFields(logrus.Fields{
		"cluster-name": clusterName,
	}).Debug("Searching for cluster in EMM configuration")

	for _, cluster := range c.Clusters {
		if cluster.Name == clusterName {
			return cluster
		}
	}

	return nil
}

// parseEMMConfig reads the EMM YAML configuration file and creates a construct with all the modules and submodules
// defined in the configuration file
func parseEMMConfig() *Config {

	logger.Debug("Reading EMM configuration file")

	configFile, err := ioutil.ReadFile(defaultEMMConfigFile)

	if err == nil {
		var emmConfig Config

		logger.Debug("Parsing the configuration file")

		err = yaml.Unmarshal(configFile, &emmConfig)

		if err != nil {

			logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("Could not parse EMM configuration file successfully")

		} else {
			return &emmConfig
		}
	}

	return nil
}
