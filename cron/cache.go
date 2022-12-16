package cron

import (
	"botApiStats/config"
	"botApiStats/dal"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"botApiStats/middleware"
	"botApiStats/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var ctx = context.Background()
var r = query.Use(dal.Db).BotResponse
var d = query.Use(dal.Db).DailyIntent

var maxN, _ = strconv.Atoi(config.Get("max_n"))

// DailyUpdateCacheCron 每天凌晨2点，整理前一天的TopN保存到daily_intent表
func DailyUpdateCacheCron() {
	c := cron.New()
	fmt.Println(time.Now(), "每天凌晨2点执行一次")
	EntryID, err := c.AddFunc("0 0 2 * * ?", DailyUpdateCacheJob)
	fmt.Println(time.Now(), EntryID, err)
	c.Start()
	fmt.Println("定时任务UpdateCacheJob已启动")
}

func DailyUpdateCacheJob() {
	fmt.Println(time.Now(), "开始更新缓存")
	CollectOneDayIntentCount(1) //整理昨日topN保存到daily_intent表

	UpdateRecentNDayTopN(1)  //缓存昨天topN结果
	UpdateRecentNDayTopN(7)  //缓存过去7天topN结果
	UpdateRecentNDayTopN(30) //缓存过去30天topN结果
	UpdateRecentNDayTopN(90) //缓存过去90天topN结果
	fmt.Println(time.Now(), "缓存更新完成")
}

// UpdateDailyIntentFrom3MonthAgo 手动触发，整理三个月内的response更新到daily_intent表
func UpdateDailyIntentFrom3MonthAgo() {
	for i := 1; i <= 90; i++ {
		CollectOneDayIntentCount(i)
	}
}

// UpdateRecentNDayTopN 从daily_intent表中查询并缓存最近n天的intent数据
func UpdateRecentNDayTopN(n int) {
	todayUnixDay := tool.UnixEpochDaysSince1970(nil)
	startDay := todayUnixDay - n
	endDay := todayUnixDay - 1

	result, err := r.SelectTopNIntentByDailyRank(int64(startDay), int64(endDay), maxN, "")

	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	beforeTomorrow := 24*time.Hour - tool.TodayUsedTimeDuration()
	cacheKey := fmt.Sprint("topn_recent", n, "day")
	err = middleware.Rdb.Set(ctx, cacheKey, data, beforeTomorrow).Err()
	fmt.Println(cacheKey, " expiration after:", beforeTomorrow)
	if err != nil {
		panic(err)
	}
}

// CollectOneDayIntentCount 整理某一天的Intent数据写入daily_intent表，入参dayAgo表示距今已有几天，昨天dayAgo=1
func CollectOneDayIntentCount(daysAgo int) []*model.DailyIntent {
	startTime := time.Now().AddDate(0, 0, -daysAgo).Add(-tool.TodayUsedTimeDuration())
	endTime := startTime.AddDate(0, 0, 1)

	var topIntents []model.IntentResult
	topIntents, err := r.SelectTopNIntentByTime(startTime.Unix(), endTime.Unix(), maxN, "")
	if err != nil {
		panic(err)
	}

	var dayIntentData []*model.DailyIntent
	for _, intent := range topIntents {
		dayIntentData = append(dayIntentData, &model.DailyIntent{
			Date:     int64(tool.UnixEpochDaysSince1970(nil) - daysAgo),
			AgentId:  intent.AgentId,
			IntentId: intent.IntentId,
			Count:    intent.Count})
	}
	err = d.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(dayIntentData, len(dayIntentData))
	if err != nil {
		panic(err)
	}

	return dayIntentData
}
