package main

import (
	"botApiStats/api"
	"botApiStats/dal"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	dal.Rdb.Set(context.Background(), fmt.Sprint(time.Now()), "redis attached", 0)

	r := gin.Default()
	//r.Use(gin.Logger()) //todo 日志改用 https://github.com/sirupsen/logrus
	r.Use(gin.Recovery())

	//对外Api
	r.POST("/chatbot/v1alpha1/agents/:agentId/channels/:channelId/getReply", api.GetBotResponse)
	r.POST("/chatbot/v1alpha1/agents/:agentId/stats/topn", api.GetTopNIntentByTimeDuration)

	//调试Api
	r.POST("/dev/topn", api.GetTopNIntentByExactTime)
	r.GET("/dev/ip", api.GetPublicIp)
	r.GET("/dev/count", api.CountAllRecord)
	r.POST("/dev/flush", api.UpdateStatsCache)
	r.GET("/dev/time", api.ShowUnixTimeInfo)
	r.GET("/dev/init", api.FillTestData)

	//api文档资源
	r.StaticFile("/v2/swagger.json", "./docs/swagger.json")
	r.StaticFile("/docs", "./templates/redoc.html")
	fmt.Println("Documentation served at http://127.0.0.1:8000/docs")

	r.Run(":8000")
}
