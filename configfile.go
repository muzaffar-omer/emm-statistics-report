package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const defaultEMMConfigFile = "emm-config.yaml"

type Config struct {
	Clusters       []Cluster       `yaml:"clusters"`
	BusinessLogics []BusinessLogic `yaml:"configurations"`
}

//######################### Main Modules ##################################
type BusinessLogic struct {
	Name          string                `yaml:"name"`
	Collectors    []string              `yaml:"collectors"`
	Distributors  []string              `yaml:"distributors"`
	LogicalServer AssignedLogicalServer `yaml:"assigned-logical-server"`
}

type Cluster struct {
	Name           string          `yaml:"name"`
	Username       string          `yaml:"username"`
	Password       string          `yaml:"password"`
	LogicalServers []LogicalServer `yaml:"logical-servers"`
}

//######################### Sub-Modules ##################################
type AssignedLogicalServer struct {
	Name    string `yaml:"name"`
	Cluster string `yaml:"cluster"`
}

type LogicalServer struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"'`
	Port     string `yaml:"port"`
}

func parseEMMConfig() *Config {
	var emmConfig Config

	configFile, err := ioutil.ReadFile(defaultEMMConfigFile)

	if err == nil {
		yaml.Unmarshal(configFile, &emmConfig)
		fmt.Printf("First cluster name : %v\n", emmConfig.Clusters[0].Name)
		fmt.Printf("First configuration is : %s \n", emmConfig.BusinessLogics[0].Name)
		fmt.Printf("Configurations collectors : %v \n", emmConfig.BusinessLogics[0].Collectors)
		fmt.Printf("Assigned Logical Server : %v \n", emmConfig.BusinessLogics[0].LogicalServer.Name)
	}

	return nil
}
