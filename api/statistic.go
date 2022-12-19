package api

import (
	"botApiStats/cache"
	"botApiStats/dal"
	"botApiStats/dal/enum"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"botApiStats/sync"
	"botApiStats/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var r = query.Use(dal.Db).BotResponse

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
	result, err := r.SelectTopNIntentByTime(req.StartTime, req.EndTime, req.N, req.AgentId)
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
// @Router			/chatbot/v1alpha1/agents/{agentId}/stats/topn [post]
func GetTopNIntentByTimeDuration(c *gin.Context) {
	agentId := c.Param("agentId")
	if agentId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "请求参数有误: agentId")
		return
	}

	var req model.DurationStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求参数有误:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "请求参数有误: "+fmt.Sprint(err))
		return
	}

	var result = make([]model.IntentResult, 0)
	//从redis获取缓存记录
	cacheValueBytes, err := dal.Rdb.Get(context.Background(), req.Duration.String()).Bytes()
	//检查Duration对应的缓存是否存在
	if err == nil && cacheValueBytes != nil && len(cacheValueBytes) > 0 {
		cacheResult := make([]model.IntentResult, 0)
		err = json.Unmarshal(cacheValueBytes, &cacheResult)
		//过滤缓存结果中非当前agentId的数据
		for i := range cacheResult {
			if cacheResult[i].AgentId == agentId {
				result = append(result, cacheResult[i])
			}
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("cache hit, req.Duration=", req.Duration)
	} else {
		//缓存未命中，从db查询
		var queryEndDay = tool.TodayUnixDay() - 1
		var queryStartDay int

		switch req.Duration {
		case enum.Yesterday:
			queryStartDay = tool.TodayUnixDay() - 1
		case enum.Recent7day:
			queryStartDay = tool.TodayUnixDay() - 7
		case enum.Recent30day:
			queryStartDay = tool.TodayUnixDay() - 30
		case enum.Recent90day:
			queryStartDay = tool.TodayUnixDay() - 90
		case enum.LastWeek:
			start, end := tool.LastWeekUnixTimeRange()
			queryStartDay, queryEndDay = tool.UnixEpochDaysSince1970(&start), tool.UnixEpochDaysSince1970(&end)
		case enum.LastMonth:
			start, end := tool.LastMonthUnixTimeRange()
			queryStartDay, queryEndDay = tool.UnixEpochDaysSince1970(&start), tool.UnixEpochDaysSince1970(&end)
		}
		result, err = r.SelectTopNIntentByDailyRank(int64(queryStartDay), int64(queryEndDay), req.N, agentId)
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "内部错误，请联系管理员")
		panic(err)
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

// @Summary		整理历史数据并刷新缓存
// @Description	整理90天内的bot_response表数据并写入daily_intent；刷新昨日、近7天、近30天、近90天的TopN缓存结果
// @Tags			开发测试
// @Accept			json
// @Produce		json
// @Success		200	body	string
// @Router			/flush [get]
func UpdateStatsCache(c *gin.Context) {
	sync.MergeLast3MonthIntentData()
	cache.UpdateAllCache()
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

// @Summary		打印当前unixTime等信息
// @Accept			json
// @Description	打印现在、今天、昨天、x天前等时间段的unix秒，unix天
// @Tags			开发测试
// @Produce		json
// @Success		200	body	string
// @Router			/time [get]
func ShowUnixTimeInfo(c *gin.Context) {
	todayStartTimeUnix := time.Now().Add(-tool.TodayUsedTimeDuration()).Unix()
	oneDayUnixTime := time.Unix(60*60*24, 0).Unix()

	wstart, wend := tool.LastWeekUnixTimeRange()
	mstart, mend := tool.LastMonthUnixTimeRange()
	result := [][4]any{
		{"now", time.Now().Unix()},
		{"yesterday_start", todayStartTimeUnix - oneDayUnixTime},
		{"yesterday_end", todayStartTimeUnix},
		{"recent7day_start", todayStartTimeUnix - 7*oneDayUnixTime},
		{"recent7day_end", todayStartTimeUnix},
		{"recent30day_start", todayStartTimeUnix - 30*oneDayUnixTime},
		{"recent30day_end", todayStartTimeUnix},
		{"recent90day_start", todayStartTimeUnix - 90*oneDayUnixTime},
		{"recent90day_end", todayStartTimeUnix},
		{"last_week_start", wstart.Unix()},
		{"last_week_end", wend.Unix()},
		{"last_month_start", mstart.Unix()},
		{"last_month_end", mend.Unix()},
	}
	for i := range result {
		result[i][0] = strings.ToUpper(result[i][0].(string))
		t := time.Unix(result[i][1].(int64), 0)
		result[i][2] = t.String()
		result[i][3] = tool.UnixEpochDaysSince1970(&t)
	}
	c.IndentedJSON(http.StatusOK, result)

}

// @Summary		写入一些模拟调用机器人api数据
// @Accept			json
// @Description	agent_id=[a1,a2] intent_id=[0,1,2,3]
// @Tags			开发测试
// @Produce		json
// @Success		200	body	string
// @Router			/init [get]
func FillTestData(c *gin.Context) {
	size := FillTestBotResponse()
	c.IndentedJSON(http.StatusOK, fmt.Sprint("insert ", size, " record to table bot_response"))
}

func FillTestBotResponse() int {
	now := time.Now()
	var simpleResponseData = []*model.BotResponse{
		//昨天
		{CreatedAt: now.Add(-24 * time.Hour).Unix(), AgentId: "a1", IntentId: "0", IntentName: "intent0_newName"},
		{CreatedAt: now.Add(-24 * time.Hour).Unix(), AgentId: "a1", IntentId: "1", IntentName: "intent1"},
		{CreatedAt: now.Add(-24 * time.Hour).Unix(), AgentId: "a2", IntentId: "2", IntentName: "intent2"},
		{CreatedAt: now.Add(-24 * time.Hour).Unix(), AgentId: "a2", IntentId: "1", IntentName: "intent1"},
		//3天前
		{CreatedAt: now.Add(-3 * 24 * time.Hour).Unix(), AgentId: "a1", IntentId: "0", IntentName: "intent0"},
		{CreatedAt: now.Add(-3 * 24 * time.Hour).Unix(), AgentId: "a1", IntentId: "2", IntentName: "intent2"},
		{CreatedAt: now.Add(-3 * 24 * time.Hour).Unix(), AgentId: "a1", IntentId: "3", IntentName: "intent3"},
		{CreatedAt: now.Add(-3 * 24 * time.Hour).Unix(), AgentId: "a2", IntentId: "2", IntentName: "intent2"},
		//10天前
		{CreatedAt: now.Add(-10 * 24 * time.Hour).Unix(), AgentId: "a2", IntentId: "2", IntentName: "intent2"},
		{CreatedAt: now.Add(-10 * 24 * time.Hour).Unix(), AgentId: "a2", IntentId: "3", IntentName: "intent3"},
		{CreatedAt: now.Add(-10 * 24 * time.Hour).Unix(), AgentId: "a2", IntentId: "3", IntentName: "intent3"},
	}
	err := r.CreateInBatches(simpleResponseData, len(simpleResponseData))
	if err != nil {
		panic(err)
		return 0
	}
	return len(simpleResponseData)
}
