package cmd

import (
	"fmt"
	"github.com/ricecake/sarge/models"
	"github.com/spf13/cobra"
	"time"
)

var year, quarter uint

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Deleivers a salary report for a specific quarter, for each department",
	Long: `Deleivers a salary report for a specific quarter, for each department

By default, this command will return a report for the current year and quarter, but
diffent years and/or quarter can be specified on the command line.
`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := models.Connect()
		if err != nil {
			panic(err)
		}
		defer db.Close()
		var deps []models.Department
		db.Find(&deps)
		for _, dep := range deps {
			s, _ := dep.GetSalaryByTimeRange(yearQuarter(year, quarter))
			fmt.Printf("%s: %0.02f\n", dep.Name, s)
		}
	},
}

var now time.Time

func init() {
	rootCmd.AddCommand(reportCmd)

	now = time.Now()
	currentQuarter := 1 + (uint(now.Month())%12)/3

	reportCmd.Flags().UintVar(&year, "year", uint(now.Year()), "Year to report on")
	reportCmd.Flags().UintVar(&quarter, "quarter", currentQuarter, "The quarter to report on")
}

func yearQuarter(year, quarter uint) (begin, end time.Time) {
	qStartMonth := 1 + 3*(quarter-1)
	qEndMonth := 1 + 3*quarter

	begin = time.Date(int(year+(qStartMonth/12)), time.Month(qStartMonth%12), 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(int(year+(qEndMonth/12)), time.Month(qEndMonth%12), 1, 0, 0, 0, 0, time.UTC)

	return begin, end
}
