package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// These variables must be passed during build
var version string
var build string
var toolName string

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)

	app := CreateCliApp()
	app.Version = fmt.Sprintf("%s - build %s", version, build)
	app.Run(os.Args)
}
