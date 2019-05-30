package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)

	//os.Args = []string{"emmstats", "--cluster=dev", "--lserver=Server11", "--group-by=month", "throughput"}
	//os.Args = []string{"emmstats", "h"}
	//os.Args = []string{"emmstats", "--stream=4GLTE_INPUT_CDRs", "--group-by=month", "throughput"}
	os.Args = []string{"emmstats", "--lserver=Server11", "--group-by=month", "--cluster=dev", "--output-file=output.csv", "--format=txt", "throughput"}

	app := CreateCliApp()
	app.Version = "1.0"
	app.Run(os.Args)
}
