package models

import (
	"math/big"
	"time"
	//"github.com/jinzhu/gorm"
	//"github.com/spf13/viper"
)

type Employee struct {
	Number    uint      `gorm:"column:emp_no"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	BirthDate time.Time `gorm:"column:birth_date"`
	Gender    string    `gorm:"column:gender"`
	HireDate  time.Time `gorm:"column:hire_date"`
}

func (Employee) TableName() string {
	return "employees"
}

func (this Employee) GetSalaryByTimeRange(from time.Time, to time.Time) (rangeSalary float64, err error) {
	startWindow := from.Format("2006-01-02")
	endWindow := to.Format("2006-01-02")

	db := GetDb()

	var salaries []Salary
	db.
		Where(&Salary{EmployeeNumber: this.Number}).
		Where(`(from_date <= ? and to_date >= ?)
		or (from_date <= ? and to_date >= ?)
		or (from_date >= ? and to_date <= ?)
		or (from_date <= ? and to_date >= ?)`,
			startWindow, startWindow,
			startWindow, endWindow,
			startWindow, endWindow,
			endWindow, endWindow,
		).
		Find(&salaries)

	total := new(big.Float)
	total.SetMode(big.AwayFromZero)
	total.SetPrec(128)

	for _, salary := range salaries {
		start := salary.StartDate
		end := salary.EndDate
		if start.Equal(from) && end.Equal(to) {
			//Easy Mode
			total.Add(total, big.NewFloat(float64(salary.Salary)))
		} else {
			var intervalStart time.Time
			var intervalStop time.Time
			if start.After(from) {
				intervalStart = start
			} else {
				intervalStart = from
			}
			if end.Before(to) {
				intervalStop = end
			} else {
				intervalStop = to
			}

			interval := intervalStop.Sub(intervalStart)
			dailyRate := big.NewFloat(salary.DailyRate())
			intervalSalary := new(big.Float).Mul(dailyRate, big.NewFloat(interval.Hours()/24))

			total.Add(total, intervalSalary)
		}
	}

	totalFloat, _ := total.Float64()
	return totalFloat, nil
}
