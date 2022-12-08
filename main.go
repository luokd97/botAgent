package main

import (
	"botApiStats/api"
	"botApiStats/cron"
	"botApiStats/middleware"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	//fmt.Println("11:", time.Now())									//result	11: 2022-12-08 23:23:40.306278 +0800 CST m=+0.108688793
	//fmt.Println("22:", time.Now().Format("2006-01-01 15:01:05.000"))	//result	22: 2022-12-12 23:12:40.306
	//fmt.Println("33:", time.Now().Format("2022-09-01 15:01:05.000"))	//result	22: 8088-09-12 23:12:40.306
	middleware.Rdb.Set(context.Background(), fmt.Sprint(time.Now()), "redis attached", 0)

	r := gin.Default()

	r.POST("/bot", api.GetBotResponse)
	r.GET("/ip", api.GetPublicIp)
	r.GET("/count", api.CountAllRecord)
	r.POST("/top", api.TopIntent)
	r.POST("/stats", api.TopIntentWithStats)
	r.POST("/flush", api.UpdateStatsCache)

	r.StaticFile("/v2/swagger.json", "./docs/swagger.json")
	r.StaticFile("/docs", "./templates/redoc.html")

	println("Documentation served at http://127.0.0.1:8000/docs")
	cron.UpdateCacheCron() //启动定时任务
	panic(r.Run(":8000"))
}
