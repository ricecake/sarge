package models

import (
	"time"
	//"github.com/jinzhu/gorm"
	//"github.com/spf13/viper"
)

type Salary struct {
	EmployeeNumber uint      `gorm:column:emp_no`
	Salary         float64   `gorm:column:salary`
	StartDate      time.Time `gorm:column:from_date`
	EndDate        time.Time `gorm:column:to_date`
}

func (Salary) TableName() string {
	return "salaries"
}

func (this Salary) Duration() uint {
	rawDuration := this.EndDate.Sub(this.StartDate)
	return uint(rawDuration.Hours() / 24)
}

func (this Salary) DailyRate() float64 {
	return this.Salary / float64(this.Duration())
}
