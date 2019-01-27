package stats

import (
	"os"

	"github.com/muzaffar-omer/emm-statistics-report/database"
	"github.com/olekukonko/tablewriter"
)

type Statistical interface {
	GetStatisticsMap() map[string]float64
}

func CalculateStats(statisticalRecords []Statistical) (map[string]float64, map[string]float64, map[string]float64, map[string]float64) {

	var sum map[string]float64 = make(map[string]float64)
	var avg map[string]float64 = make(map[string]float64)
	var min map[string]float64 = make(map[string]float64)
	var max map[string]float64 = make(map[string]float64)

	var numberOfRecords int = len(statisticalRecords)

	for index, record := range statisticalRecords {

		for key, value := range record.GetStatisticsMap() {
			sum[key] += float64(value)
			avg[key] += float64(value)

			if min[key] > float64(value) {
				min[key] = float64(value)
			}

			if max[key] < float64(value) {
				max[key] = float64(value)
			}

			// If this is the last record, calculate the average
			if index+1 == numberOfRecords {
				avg[key] = avg[key] / float64(numberOfRecords)
			}
		}
	}

	return sum, avg, min, max
}

func CreateStatisticsTable(statisticalRecords []Statistical) *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	var sum map[string]float64
	var avg map[string]float64
	var min map[string]float64
	var max map[string]float64

	sum, avg, min, max = CalculateStats(statisticalRecords)

	table = tablewriter.NewWriter(os.Stdout)

	// Create header
	header := make([]string, 0, len(sum)+1)
	header = append(header, "")

	for key, _ := range sum {
		header = append(header, key)
	}
	table.SetHeader(header)

	var statsRow []string

	// Fill table data
	for rowNum := 0; rowNum < 4; rowNum += 1 {

		// Fill summations row
		if rowNum == 0 {
			statsRow = make([]string, len(header))
			statsRow[0] = "Sum"

			for colNum := 1; colNum < len(header); colNum += 1 {
				statsRow[colNum] = database.GetFormattedNumber(sum[header[colNum]])
			}

			table.Append(statsRow)
		} else if rowNum == 1 { // Fill average row
			statsRow = make([]string, len(header))
			statsRow[0] = "Avg"

			for colNum := 1; colNum < len(header); colNum += 1 {
				statsRow[colNum] = database.GetFormattedNumber(avg[header[colNum]])
			}

			table.Append(statsRow)
		} else if rowNum == 2 { // Fill min row
			statsRow = make([]string, len(header))
			statsRow[0] = "Min"

			for colNum := 1; colNum < len(header); colNum += 1 {
				statsRow[colNum] = database.GetFormattedNumber(min[header[colNum]])
			}

			table.Append(statsRow)
		} else if rowNum == 3 { // Fill max row
			statsRow = make([]string, len(header))
			statsRow[0] = "Max"

			for colNum := 1; colNum < len(header); colNum += 1 {
				statsRow[colNum] = database.GetFormattedNumber(max[header[colNum]])
			}

			table.Append(statsRow)
		}
	}

	return table
}
