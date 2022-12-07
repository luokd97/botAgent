package api

import (
	"botApiStats/cron"
	"botApiStats/dal"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var r = query.Use(dal.Db).BotResponse

func AddRecord(id string, name string) {
	r.Create(&model.BotResponse{IntentId: id, IntentName: name})
}

// @Summary 获取知识点TopN（无缓存）
// @Description 点击200 Successful Response查看具体接口返回格式
// @Tags botAgent
// @Accept json
// @Produce json
// @Param   n	body	model.StatsRequest	true  " "
// @Success 200 {object} []model.StatsResponse
// @Router /top [post]
func TopIntent(c *gin.Context) {
	n := 3
	startTime := "1950-01-01 00:00:00.00"
	endTime := "2099-01-01 00:00:00.00"
	var req model.StatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("request param illegal", err)
	} else {
		if req.N != 0 {
			n = req.N
		}
		if req.StartTime != "" {
			startTime = req.StartTime
		}
		if req.EndTime != "" {
			endTime = req.EndTime
		}
	}

	result, err := r.SelectTopIntentIdByTime(startTime, endTime, n)
	if err != nil {
		panic(err)
	}

	//展示最新intent_name
	for _, v := range result {
		intentId := v["intent_id"]
		fmt.Println("intentId=", intentId)
		intentName, _ := r.GetIntentNameByIntentId(fmt.Sprint(intentId))
		v["intent_name"] = intentName
	}
	fmt.Printf("SelectTopIntentId(%v) =%v\n", n, result)
	c.IndentedJSON(http.StatusOK, result)
}

// @Summary 获取知识点TopN
// @Description 点击200 Successful Response查看具体接口返回格式
// @Tags botAgent
// @Accept json
// @Produce json
// @Param   n	body	model.StatsRequest	true  " "
// @Success 200 {object} []model.StatsResponse
// @Router /stats [post]
func TopIntentWithStats(c *gin.Context) {

}

// @Summary 刷新统计数据缓存
// @Description 点击200 Successful Response查看具体接口返回格式
// @Tags botAgent
// @Accept json
// @Produce json
// @Success 200 body string
// @Router /flush [post]
func UpdateStatsCache(c *gin.Context) {
	cron.UpdateCacheJob()
	c.IndentedJSON(http.StatusOK, "执行刷新缓存任务")
}

func CountAllRecord(c *gin.Context) {
	count, err := r.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	c.IndentedJSON(http.StatusOK, count)
}

func FillSampleBotResponse() {
	var records = []model.BotResponse{
		{IntentId: "0", IntentName: "intent0"},
		{IntentId: "1", IntentName: "intent1"},
		{IntentId: "2", IntentName: "intent2"},
		{IntentId: "3", IntentName: "intent3"},
		{IntentId: "4", IntentName: "intent4"},
	}
	dal.Db.CreateInBatches(&records, len(records))
}
