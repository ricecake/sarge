package models

import (
	"time"
	//"github.com/jinzhu/gorm"
	//"github.com/spf13/viper"
)

type Assignment struct {
	EmployeeNumber   uint      `gorm:"column:emp_no"`
	DepartmentNumber string    `gorm:"column:dept_no"`
	StartDate        time.Time `gorm:"column:from_date"`
	EndDate          time.Time `gorm:"column:to_date"`
}

func (Assignment) TableName() string {
	return "dept_emp"
}

func (this Assignment) Employee() Employee {
	db := GetDb()

	var employee Employee
	db.Where(&Employee{Number: this.EmployeeNumber}).Find(&employee)
	return employee
}

func (this Assignment) Department() Department {
	db := GetDb()

	var department Department
	db.Where(&Department{Number: this.DepartmentNumber}).Find(&department)
	return department
}
