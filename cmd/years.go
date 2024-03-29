package cmd

import "time"

const (
	startYear = 2019
)

func allYears() []int {
	currentYear := time.Now().Year()
	finalYear := currentYear + 10
	years := pastYears()

	for i := currentYear + 1; i <= finalYear; i++ {
		years = append(years, i)
	}

	return years
}

func pastYears() []int {
	currentYear := time.Now().Year()
	years := make([]int, 0)

	for i := startYear; i <= currentYear; i++ {
		years = append(years, i)
	}

	return years
}
