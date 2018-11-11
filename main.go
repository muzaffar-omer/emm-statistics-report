package main

import (
	"emm-statistics-report/configuration"
	"fmt"
)

func main() {
	config := configuration.Config{}
	config.Init()

	fmt.Printf("%s", config)
}
