package api

import (
	"botApiStats/cron"
	"botApiStats/dal"
	"botApiStats/dal/enum"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"botApiStats/middleware"
	"botApiStats/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var r = query.Use(dal.Db).BotResponse

func AddRecord(id string, name string) {
	r.Create(&model.BotResponse{IntentId: id, IntentName: name})
}

// @Summary		按精确时间获取知识点TopN
// @Description	点击200 Successful Response查看具体接口返回格式
// @Tags			开发测试
// @Accept			json
// @Produce		json
// @Param			n	body		model.StatsRequest	true	" "
// @Success		200	{object}	[]model.StatsResponse
// @Router			/top [post]
func TopIntent(c *gin.Context) {
	n := 3
	startTime := 0
	endTime := 2147483647 //At 03:14:08 UTC on Tuesday, 19 January 2038
	var req model.StatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("request param illegal", err)
	} else {
		if req.N != 0 {
			n = req.N
		}
		if req.StartTime != 0 {
			startTime = req.StartTime
		}
		if req.EndTime != 0 {
			endTime = req.EndTime
		}
	}

	resultMap, err := r.SelectTopIntentIdByTime(int64(startTime), int64(endTime), n)
	if err != nil {
		panic(err)
	}
	var result []model.StatsResponse
	tool.Map2struct(resultMap, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("SelectTopIntentId(%v) =%v\n", n, result)
	c.IndentedJSON(http.StatusOK, result)
}

// @Summary		获取知识点TopN
// @Description	点击200 Successful Response查看具体接口返回格式
// @Tags			botAgent
// @Accept			json
// @Produce		json
// @Param			n	body		model.CommonIntentRequest	true "控制参数"
// @Success		200	{object}	[]model.StatsResponse
// @Router			/stats [post]
func TopIntentWithStats(c *gin.Context) {
	n := 3
	duration := enum.Yesterday
	var req model.CommonIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("request param illegal", err)
	} else {
		if req.N != 0 {
			n = req.N
		}
		if &req.Duration != nil {
			duration = req.Duration
		}
	}

	var startUnixTime int64
	var endUnixTime int64

	now := time.Now()
	todayTime := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second
	endUnixTime = now.Add(-todayTime).Unix()

	var err error
	var result []model.StatsResponse
	switch duration {
	case enum.Yesterday:
		startUnixTime = now.Add(-1 * 24 * time.Hour).Add(-todayTime).Unix()

		//缓存检查
		s, err := middleware.Rdb.Get(context.Background(), "topn_yesterday").Bytes()
		if err == nil {
			json.Unmarshal(s, &result)
			fmt.Println("yesterday cache hit")
		}
	case enum.Recent7day:
		startUnixTime = now.Add(-7 * 24 * time.Hour).Add(-todayTime).Unix()

		//缓存检查
		s, err := middleware.Rdb.Get(context.Background(), "topn_recent7day").Bytes()
		if err == nil {
			json.Unmarshal(s, &result)
			fmt.Println("Recent7day cache hit")
		}
	case enum.Recent30day:
		startUnixTime = now.Add(-30 * 24 * time.Hour).Add(-todayTime).Unix()

		//缓存检查
		s, err := middleware.Rdb.Get(context.Background(), "topn_recent30day").Bytes()
		if err == nil {
			json.Unmarshal(s, &result)
			fmt.Println("Recent30day cache hit")
		}
	case enum.Recent90day:
		startUnixTime = now.Add(-90 * 24 * time.Hour).Add(-todayTime).Unix()

		//缓存检查
		s, err := middleware.Rdb.Get(context.Background(), "topn_recent90day").Bytes()
		if err == nil {
			json.Unmarshal(s, &result)
			fmt.Println("Recent90day cache hit")
		}
	case enum.LastWeek:
		daysUtilLastSunday := now.Weekday()
		if now.Weekday() == time.Sunday {
			daysUtilLastSunday = 7
		}
		startUnixTime = now.Add(time.Duration(-24*(daysUtilLastSunday+7)) * time.Hour).Unix()
		endUnixTime = now.Add(time.Duration(-24*daysUtilLastSunday) * time.Hour).Unix()
	case enum.LastMonth:
		startUnixTime = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, nil).Unix()
		endUnixTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, nil).Unix()
	}

	//缓存未命中，从db查询
	if result == nil {
		resultMap, _ := r.SelectTopIntentByDailyRank(startUnixTime, endUnixTime, n)
		tool.Map2struct(resultMap, &result)
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, result)

}

// @Summary		手动刷新统计数据缓存
// @Description	刷新昨日、近7天、近30天、近90天的TopN缓存结果
// @Tags			开发测试
// @Accept			json
// @Produce		json
// @Success		200	body	string
// @Router			/flush [post]
func UpdateStatsCache(c *gin.Context) {
	cron.UpdateCacheJob()
	c.IndentedJSON(http.StatusOK, "现在开始执行刷新缓存任务")
}

func CountAllRecord(c *gin.Context) {
	count, err := r.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	c.IndentedJSON(http.StatusOK, count)
}
