package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)

	app := CreateCliApp()
	app.Version = "1.0"
	app.Run(os.Args)
}
