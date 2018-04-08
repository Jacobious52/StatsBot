package model

import "sort"

type UniformUserExporter struct {
	Stats Timeseries
}

func (u UniformUserExporter) Export() (xVals []float64, yVals map[string][]float64) {
	yVals = make(map[string][]float64)
	users := make(map[string]struct{})

	var sortedDays []int
	for day := range u.Stats {
		sortedDays = append(sortedDays, day)
	}
	sort.Ints(sortedDays)

	for _, day := range sortedDays {
		freq := u.Stats[day]
		xVals = append(xVals, float64(day))
		for user := range freq {
			users[user] = struct{}{}
		}
	}

	for _, day := range sortedDays {
		freq := u.Stats[day]
		for user := range users {
			if _, ok := freq[user]; !ok {
				yVals[user] = append(yVals[user], 0.0)
				continue
			}
			yVals[user] = append(yVals[user], float64(freq[user]))
		}
	}

	return
}
