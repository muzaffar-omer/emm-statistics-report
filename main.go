package main

import (
	config "emm-statistics-report/configuration"
	"emm-statistics-report/database"
	"emm-statistics-report/stats"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	logger := config.Log()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true})

	switch config.CmdConfig.OperationType() {
	case 1:
		OperationGroupedProcessedInOut()
		break
	default:
		OperationGroupedProcessedInOut()
	}

	fmt.Println("Test CI/CD 1 ...")
}

// Possible operations:
// 1 - Get processed input/output grouped by minute, hour, day, or month
func OperationGroupedProcessedInOut() {
	var totalGroupedProcessedInOut database.TotalGroupedProcessedInOut
	var statisticalRecords []stats.Statistical

	table := tablewriter.NewWriter(os.Stdout)

	if stream := config.GetStreamInfo(config.CmdConfig.Stream()); stream != nil {

		rows := database.GetGroupedStreamProcessedInOut(stream, config.CmdConfig.GroupBy())

		if rows != nil {

			for rows.Next() {
				totalGroupedProcessedInOut = database.TotalGroupedProcessedInOut{}
				rows.StructScan(&totalGroupedProcessedInOut)
				statisticalRecords = append(statisticalRecords, totalGroupedProcessedInOut)
				table.Append(totalGroupedProcessedInOut.AsArray())
			}

			table.SetHeader(totalGroupedProcessedInOut.Header())
			table.Render()

			var sum map[string]float64
			var avg map[string]float64
			var min map[string]float64
			var max map[string]float64

			sum, avg, min, max = stats.CalculateStats(statisticalRecords)

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

			table.Render()
		}
	}
}
