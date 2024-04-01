package cmd

import (
	"fmt"
	"io"
)

var (
	borrowedSum                      float64
	borrowedYear                     int
	repaymentStartYear               int
	sharingPercentage                float64
	startingSalary                   float64
	expectedSalaryIncreasePercentage float64
	expectedCPIIncreasePercentage    float64
	repaymentYears                   int
)

// thresholdValue returns the threshold value for the given year based on the CPI rates. The threshold value is the
// value of the borrowed sum times two compounded with CPI rates for the years between the borrowed year and the given year.
func thresholdValue(year int) float64 {
	threshold := borrowedSum * 2
	cpir := cpiRates(year)
	for i := borrowedYear; i <= year; i++ {
		threshold *= 1 + cpir[i]/100
	}

	return threshold
}

// cpiRates returns the map of CPI rates from 2019 to given year. The historical CPI rates are hardcoded. The future
// rates are filled in with the expectedCPIIncreasePercentage.
func cpiRates(year int) map[int]float64 {
	// Historical CPI rates from https://www-genesis.destatis.de/genesis/online?operation=abruftabelleBearbeiten&levelindex=1&levelid=1712007485333&auswahloperation=abruftabelleAuspraegungAuswaehlen&auswahlverzeichnis=ordnungsstruktur&auswahlziel=werteabruf&code=61111-0001&auswahltext=&werteabruf=Value+retrieval#abreadcrumb
	cpir := map[int]float64{
		2019: 1.4,
		2020: 0.5,
		2021: 3.1,
		2022: 6.9,
		2023: 5.9,
	}

	for i := 2024; i <= year; i++ {
		cpir[i] = expectedCPIIncreasePercentage
	}

	return cpir
}

// salary returns the salary for the given year based on the starting salary and the expected salary increase percentage.
func salary(year int) float64 {
	s := startingSalary
	for i := repaymentStartYear; i < year; i++ {
		s *= 1 + expectedSalaryIncreasePercentage/100
	}

	return s
}

// repayment returns the payment for the given year based on the salary, sharing percentage, and threshold value.
func repayment(year int) float64 {
	return salary(year) * sharingPercentage / 100
}

// paymentSchedule prints the payment schedule for the given year range to a io.Writer.
func paymentSchedule(w io.Writer) {
	totalRepaid := 0.0
	fmt.Fprintf(w, "Year\tSalary\tRepayment\tThreshold\n")
	for i := repaymentStartYear; i <= repaymentStartYear+repaymentYears; i++ {
		repaid := repayment(i)
		totalRepaid += repaid
		if totalRepaid > thresholdValue(i) {
			delta := totalRepaid - thresholdValue(i)
			totalRepaid = thresholdValue(i)

			fmt.Fprintf(w, "%d\t%.2f\t%.2f\t%.2f\n", i, salary(i), repaid-delta, thresholdValue(i))
			break
		}
		fmt.Fprintf(w, "%d\t%.2f\t%.2f\t%.2f\n", i, salary(i), repayment(i), thresholdValue(i))
	}
	fmt.Fprintf(w, "\nTotal repaid: %.2f\n", totalRepaid)
	repaidFraction := totalRepaid / borrowedSum
	fmt.Fprintf(w, "Multiples of borrowed sum repaid: %.2f\n", repaidFraction)
}

func printVariables(w io.Writer) {
	fmt.Fprintf(w, "Here are the values you entered:\n")
	fmt.Fprintf(w, "Borrowed sum: %.2f\n", borrowedSum)
	fmt.Fprintf(w, "Borrowed year: %d\n", borrowedYear)
	fmt.Fprintf(w, "Repayment start year: %d\n", repaymentStartYear)
	fmt.Fprintf(w, "Sharing percentage: %.2f\n", sharingPercentage)
	fmt.Fprintf(w, "Starting salary: %.2f\n", startingSalary)
	fmt.Fprintf(w, "Expected salary increase percentage: %.2f\n", expectedSalaryIncreasePercentage)
	fmt.Fprintf(w, "Expected CPI increase percentage: %.2f\n", expectedCPIIncreasePercentage)
	fmt.Fprintf(w, "Repayment years: %d\n", repaymentYears)
}
