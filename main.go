package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)

	emmConfig = parseEMMConfig()

	//os.Args = []string{"emmstats", "--cluster=dev", "--lserver=Server11", "--group-by=month", "throughput"}
	//os.Args = []string{"emmstats", "h"}
	os.Args = []string{"emmstats", "--stream=4GLTE_INPUT_CDRs", "--group-by=month", "throughput"}

	app := CreateCliApp()
	app.Run(os.Args)
}
