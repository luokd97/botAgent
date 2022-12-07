package main

import (
	"botApiStats/api"
	"github.com/gin-gonic/gin"
)

func main() {
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
	panic(r.Run(":8000"))
}
