package main

import (
	"botApiStats/dal"
	"botApiStats/dal/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "dal/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(dal.Db) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.BotResponse{})

	// Generate Type Safe API with Dynamic SQL defined on Method interface for `model.User` and `model.Company`
	g.ApplyInterface(func(model.Method) {}, model.BotResponse{})

	// Generate the code
	g.Execute()
}
