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
func paymentSchedule(w io.Writer, styles *Styles) {
	title := styles.Highlight.Render("📊 Here's your payment schedule")
	headerYear := styles.HeaderText.Render("📅 Year")
	headerSalary := styles.HeaderText.Render("💶 Salary")
	headerRepayment := styles.HeaderText.Render("💸 Repayment")
	headerThreshold := styles.HeaderText.Render("💰 Threshold Sum")
	fmt.Fprintf(w, "%s\n\n", title)
	fmt.Fprintf(w, "%s%s%s%s\n\n", headerYear, headerSalary, headerRepayment, headerThreshold)

	totalRepaid := 0.0
	totalYears := 0
	finalYear := repaymentStartYear + repaymentYears
	for i := repaymentStartYear; i <= repaymentStartYear+repaymentYears; i++ {
		totalYears += 1
		repaid := repayment(i)
		totalRepaid += repaid

		valueYear := styles.TableText.Render(fmt.Sprintf("%d", i))
		valueSalary := styles.TableText.Render(fmt.Sprintf("%.2f", salary(i)))
		valueRepayment := styles.TableText.Render(fmt.Sprintf("%.2f", repayment(i)))
		valueThreshold := styles.TableText.Render(fmt.Sprintf("%.2f", thresholdValue(i)))
		if totalRepaid > thresholdValue(i) {
			finalYear = i
			delta := totalRepaid - thresholdValue(i)
			totalRepaid = thresholdValue(i)

			valueRepayment = styles.TableText.Render(fmt.Sprintf("%.2f", repaid-delta))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", valueYear, valueSalary, valueRepayment, valueThreshold)
			break
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", valueYear, valueSalary, valueRepayment, valueThreshold)
	}

	totalRepaidStr := styles.Highlight.Render("Total repaid:")
	fmt.Fprintf(w, "\n%s %.2f\n", totalRepaidStr, totalRepaid)

	repaidFraction := totalRepaid / borrowedSum
	multiplesOfBorrowedSumStr := styles.Highlight.Render("Multiples of borrowed sum repaid:")
	fmt.Fprintf(w, "%s %.2f\n", multiplesOfBorrowedSumStr, repaidFraction)

	equivalentToLoanStr := styles.Highlight.Render("Equivalent to a education loan with interest rate of")
	fmt.Fprintf(w, "%s %.2f%%", equivalentToLoanStr, equivalentLoanInterestRate(finalYear, totalYears)*100)
}

func printVariables(w io.Writer) {
	fmt.Fprintf(w, "\tBorrowed sum: %.2f\n", borrowedSum)
	fmt.Fprintf(w, "\tBorrowed year: %d\n", borrowedYear)
	fmt.Fprintf(w, "\tRepayment start year: %d\n", repaymentStartYear)
	fmt.Fprintf(w, "\tSharing percentage: %.2f\n", sharingPercentage)
	fmt.Fprintf(w, "\tStarting salary: %.2f\n", startingSalary)
	fmt.Fprintf(w, "\tExpected salary increase percentage: %.2f\n", expectedSalaryIncreasePercentage)
	fmt.Fprintf(w, "\tExpected CPI increase percentage: %.2f\n", expectedCPIIncreasePercentage)
	fmt.Fprintf(w, "\tRepayment years: %d\n", repaymentYears)
}

func equivalentLoanInterestRate(finalYear int, totalYears int) float64 {
	finalPaid := thresholdValue(finalYear)
	interest := finalPaid - borrowedSum
	years := float64(totalYears)
	if interest == 0 {
		return 0.0
	}

	return interest / borrowedSum / years
}
