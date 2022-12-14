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
	r.Create(&model.BotResponse{IntentId: id, IntentName: name, CreatedAt: time.Now().Unix()})
}

// @Summary		按精确范围统计知识点TopN
// @Description	点击200 Successful Response查看具体接口返回格式
// @Tags			开发测试
// @Accept			json
// @Produce		json
// @Param			n	body		model.ExactStatsRequest	true	" "
// @Success		200	{object}	[]model.IntentResult
// @Router			/topn [post]
func GetTopNIntentByExactTime(c *gin.Context) {
	var req model.ExactStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求参数有误:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "请求参数有误: "+fmt.Sprint(err))
		return
	}
	var result []model.IntentResult
	result, err := r.SelectTopNIntentByTime(req.StartTime, req.EndTime, req.N)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "内部错误，请联系管理员")
		panic(err)
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

// @Summary		按枚举范围统计知识点TopN
// @Description	点击200 Successful Response查看具体接口返回格式
// @Tags			botAgent
// @Accept			json
// @Produce		json
// @Param			n	body		model.DurationStatsRequest	true "控制参数"
// @Success		200	{object}	[]model.IntentResult
// @Router			/stats [post]
func GetTopNIntentByTimeDuration(c *gin.Context) {
	var req model.DurationStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求参数有误:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "请求参数有误: "+fmt.Sprint(err))
		return
	}
	var result = make([]model.IntentResult, 0)

	//检查Duration对应的缓存是否存在
	cacheHit := false
	var cacheValueBytes []byte
	var err error
	switch req.Duration {
	case enum.Yesterday:
		cacheValueBytes, err = middleware.Rdb.Get(context.Background(), "topn_yesterday").Bytes()
	case enum.Recent7day:
		cacheValueBytes, err = middleware.Rdb.Get(context.Background(), "topn_recent7day").Bytes()
	case enum.Recent30day:
		cacheValueBytes, err = middleware.Rdb.Get(context.Background(), "topn_recent30day").Bytes()
	case enum.Recent90day:
		cacheValueBytes, err = middleware.Rdb.Get(context.Background(), "topn_recent90day").Bytes()
	case enum.LastWeek:
		cacheValueBytes = nil
	case enum.LastMonth:
		cacheValueBytes = nil
	}
	if err == nil && cacheValueBytes != nil && len(cacheValueBytes) > 0 {
		err = json.Unmarshal(cacheValueBytes, &result)
		if err != nil {
			panic(err)
		}
		cacheHit = true
		fmt.Println("cache hit, req.Duration=", req.Duration)
	}

	//缓存未命中，从db查询
	if !cacheHit {
		nowTime := time.Now()
		//todayStartTime对应今天的0点0分
		todayStartTime := nowTime.Add(-tool.TodayUsedTimeDuration())

		var queryStartTime time.Time
		var queryEndTime = todayStartTime
		switch req.Duration {
		case enum.Yesterday:
			queryStartTime = todayStartTime.Add(-1 * 24 * time.Hour)
		case enum.Recent7day:
			queryStartTime = todayStartTime.Add(-7 * 24 * time.Hour)
		case enum.Recent30day:
			queryStartTime = todayStartTime.Add(-30 * 24 * time.Hour)
		case enum.Recent90day:
			queryStartTime = todayStartTime.Add(-90 * 24 * time.Hour)
		case enum.LastWeek:
			daysUtilLastSunday := nowTime.Weekday()
			if nowTime.Weekday() == time.Sunday {
				daysUtilLastSunday = 7
			}
			queryStartTime = todayStartTime.Add(time.Duration(-24*(daysUtilLastSunday+7)) * time.Hour)
			queryEndTime = todayStartTime.Add(time.Duration(-24*daysUtilLastSunday) * time.Hour)
		case enum.LastMonth:
			if nowTime.Month() == 1 {
				queryStartTime = time.Date(nowTime.Year()-1, 12, 1, 0, 0, 0, 0, time.Local)
			} else {
				queryStartTime = time.Date(nowTime.Year(), nowTime.Month()-1, 1, 0, 0, 0, 0, time.Local)
			}
			queryEndTime = queryStartTime.AddDate(0, 1, 0)
		}

		if result == nil || len(result) == 0 {
			result, err = r.SelectTopNIntentByDailyRank(queryStartTime.Unix(), queryEndTime.Unix(), req.N)
		}
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "内部错误，请联系管理员")
		panic(err)
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

// @Summary		手动刷新统计数据缓存
// @Description	刷新昨日、近7天、近30天、近90天的TopN缓存结果
// @Tags			开发测试
// @Accept			json
// @Produce		json
// @Success		200	body	string
// @Router			/flush [get]
func UpdateStatsCache(c *gin.Context) {
	cron.UpdateCacheJob()
	c.IndentedJSON(http.StatusOK, "现在开始执行刷新缓存任务")
}

// @Summary		总记录数
// @Accept			json
// @Description	返回bot_response总行数
// @Tags			开发测试
// @Produce		json
// @Success		200	body	string
// @Router			/count [get]
func CountAllRecord(c *gin.Context) {
	count, err := r.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	c.IndentedJSON(http.StatusOK, count)
}
