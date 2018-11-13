package main

import (
	"emm-statistics-report/configuration"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true})

	config := configuration.Config{}
	config.Init()

	log.WithFields(log.Fields{
		"ip":              config.Ip(),
		"port":            config.Port(),
		"password":        config.Password(),
		"username":        config.Username(),
		"logical_servers": config.Clusters(),
		"streams":         config.Streams(),
	}).Info("Configuration")
}
