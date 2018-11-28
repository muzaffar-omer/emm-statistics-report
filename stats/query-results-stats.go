package stats

type Statistical interface {
	GetStatisticsMap() map[string]int
}

func CalculateStats(statisticalRecords []Statistical) (map[string]int, map[string]float64, map[string]int, map[string]int) {

	var sum map[string]int = make(map[string]int)
	var avg map[string]float64 = make(map[string]float64)
	var min map[string]int = make(map[string]int)
	var max map[string]int = make(map[string]int)

	var numberOfRecords int = len(statisticalRecords)

	for index, record := range statisticalRecords {

		for key, value := range record.GetStatisticsMap() {
			sum[key] += value
			avg[key] += float64(value)

			if min[key] > value {
				min[key] = value
			}

			if max[key] < value {
				max[key] = value
			}

			// If this is the last record, calculate the average
			if index+1 == numberOfRecords {
				avg[key] = avg[key] / float64(numberOfRecords)
			}
		}
	}

	//for key, value := range sum {
	//	fmt.Printf("Sum %s = %d\n", key, value)
	//}

	return sum, avg, min, max
}
