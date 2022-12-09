package dal

import (
	"botApiStats/config"
	"botApiStats/dal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var Db *gorm.DB

func init() {
	fmt.Println("mysql.go init()")

	dbuser := config.Get("mysql_user")
	dbpass := config.Get("mysql_pass")
	url := config.Get("mysql_url")
	dbname := config.Get("mysql_dbname")

	dsn := dbuser + ":" + dbpass + url + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	retry := 0
	for retry < 5 && (Db == nil || err != nil) {
		Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if Db != nil && err == nil {
			break
		}
		fmt.Printf("数据库连接失败，3秒后重试，当前重试次数:%v Db==nil:%v Db.Error!=nil:%v", retry, Db == nil, Db.Error != nil)
		time.Sleep(3 * time.Second)
		retry++
	}

	Db.AutoMigrate(&model.BotResponse{}, &model.DailyIntent{})
}
