package stats

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
