package main

import (
	"encoding/csv"
	"github.com/olekukonko/tablewriter"
	"os"
)

type Report struct {
	Data [][]string
	Name string
	OutputFile string
}

func writeToFile(data [][]string, filename string, format string) {

	switch format {
	case txtFileFormat: writeTxt(data, filename); break
	case csvFileFormat: writeCSV(data, filename); break
	case xlsFileFormat: writeXLS(data, filename); break
	default: writeTxt(data, filename); break
	}
}

func writeTxt(data [][]string, filename string) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		logger.Error(err)
		return
	}

	table := tablewriter.NewWriter(file)
	table.SetHeader(data[0])
	table.AppendBulk(data[1:])
	table.Render()
}

func writeCSV(data [][]string, filename string) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		logger.Error(err)
		return
	}

	w := csv.NewWriter(file)

	for _, record := range data {
		w.Write(record)
	}

	w.Flush()
}

func writeXLS(data [][]string, filename string) {

}
