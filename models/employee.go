package models

import (
//"github.com/jinzhu/gorm"
//"github.com/spf13/viper"
)

type Employee struct {
	Number    uint      `gorm:column:emp_no`
	FirstName string    `gorm:column:first_name`
	LastName  string    `gorm:column:last_name`
	BirthDate time.Time `gorm:column:birth_date`
	Gender    string    `gorm:column:gender`
	HireDate  time.Time `gorm:column:hire_date`
}

func (Employee) TableName() string {
	return "employees"
}
