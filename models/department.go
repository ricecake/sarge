package models

import (
	"time"
	//"github.com/spf13/viper"
)

type Department struct {
	Name   string `gorm:"column:dept_name"`
	Number string `gorm:"column:dept_no"`
}

func (Department) TableName() string {
	return "departments"
}

func (this Department) GetEmployeeByTimeRange(from time.Time, to time.Time) (employees []Employee, err error) {
	startWindow := from.Format("2006-01-02")
	endWindow := to.Format("2006-01-02")

	db := GetDb()

	var assignments []Assignment
	db.
		Where(&Assignment{DepartmentNumber: this.Number}).
		Where(`(from_date <= ? and to_date >= ?)
		or (from_date <= ? and to_date >= ?)
		or (from_date >= ? and to_date <= ?)
		or (from_date <= ? and to_date >= ?)`,
			startWindow, startWindow,
			startWindow, endWindow,
			startWindow, endWindow,
			endWindow, endWindow,
		).
		Find(&assignments)

	for _, assignment := range assignments {
		employees = append(employees, assignment.Employee())
	}

	return employees, nil
}
