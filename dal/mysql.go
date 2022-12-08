package dal

import (
	"botApiStats/dal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dbuser = "root"
const dbpass = "12345678"

const url = "@tcp(mysql:3306)/"

// const url = "@tcp(localhost:3306)/"
const dbname = "bot_agent"

var Db *gorm.DB

func init() {
	dsn := dbuser + ":" + dbpass + url + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		fmt.Println(err)
	}
	Db.AutoMigrate(&model.BotResponse{}, &model.DailyIntent{})
}
