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
	middleware.Rdb.Set(context.Background(), fmt.Sprint(time.Now()), "redis attached", 0)

	r := gin.Default()
	//r.Use(gin.Logger()) //todo 日志改用 https://github.com/sirupsen/logrus
	r.Use(gin.Recovery())

	//对外Api
	r.POST("/chatbot/v1alpha1/agents/:agentId/channels/:channelId/getReply", api.GetBotResponse)
	r.POST("/topn", api.GetTopNIntentByTimeDuration)

	//调试Api
	r.POST("/v2/topn", api.GetTopNIntentByExactTime)
	r.GET("/ip", api.GetPublicIp)
	r.GET("/count", api.CountAllRecord)
	r.POST("/flush", api.UpdateStatsCache)

	//api文档资源
	r.StaticFile("/v2/swagger.json", "./docs/swagger.json")
	r.StaticFile("/docs", "./templates/redoc.html")
	fmt.Println("Documentation served at http://127.0.0.1:8000/docs")

	//启动定时任务
	cron.UpdateCacheCron()

	r.Run(":8000")
}
