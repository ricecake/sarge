package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var db *gorm.DB

func Connect() (*gorm.DB, error) {
	var err error
	if db == nil {
		username := viper.GetString("db.username")
		password := viper.GetString("db.password")
		host := viper.GetString("db.host")
		database := viper.GetString("db.database")
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, database))
	}

	return db, err
}

func SetDb(newDb *gorm.DB) {
	db = newDb
}

func GetDb() *gorm.DB {
	if db == nil {
		if _, err := Connect(); err != nil {
			panic(err)
		}
	}
	return db
}
