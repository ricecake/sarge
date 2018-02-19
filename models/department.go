package models

import (
//"github.com/jinzhu/gorm"
//"github.com/spf13/viper"
)

type Department struct {
	Name   string `gorm:"column:dept_name"`
	Number string `gorm:"column:dept_no"`
}

func (Department) TableName() string {
	return "departments"
}
