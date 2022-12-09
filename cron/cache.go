package cron

import (
	"botApiStats/dal"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"botApiStats/middleware"
	"botApiStats/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

var ctx = context.Background()
var r = query.Use(dal.Db).BotResponse
var d = query.Use(dal.Db).DailyIntent

//var datetimeFormat = "2000-01-01 00:00:00.00"

func TestCron() {
	fmt.Println("start cron test")
	c := cron.New()
	EntryID, err := c.AddFunc("@every 5s", func() { fmt.Println("Every 5s") })
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
}

func UpdateCacheCron() {
	c := cron.New()
	fmt.Println(time.Now(), "每天凌晨2点执行一次")
	EntryID, err := c.AddFunc("0 0 2 * * ?", UpdateCacheJob)
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
}

func UpdateCacheJob() {
	//更新缓存
	fmt.Println(time.Now(), "开始更新缓存")
	updateYesterdayTopN() //缓存昨天topN结果

	updateRecentNDayTopN(7)  //缓存过去7天topN结果
	updateRecentNDayTopN(30) //缓存过去30天topN结果
	updateRecentNDayTopN(90) //缓存过去90天topN结果
	fmt.Println(time.Now(), "缓存更新完成")
}

// 统计昨日topN，保存到daily_intent表
func updateYesterdayTopN() {
	dayIntents := CollectOneDayIntentCount(1)

	idNameMap, _ := r.SelectIntentIdMapIntentName()

	var yesterdayStats []model.StatsResponse
	for i := range dayIntents {
		yesterdayStats = append(yesterdayStats, model.StatsResponse{
			Cnt:        dayIntents[i].Count,
			IntentId:   dayIntents[i].IntentId,
			IntentName: idNameMap[dayIntents[i].IntentId],
		})
	}
	data, err := json.Marshal(yesterdayStats)
	if err != nil {
		panic(err)
	}
	err = middleware.Rdb.Set(ctx, "topn_yesterday", data, 48*time.Hour).Err()
	if err != nil {
		panic(err)
	}

}
func updateRecentNDayTopN(n int) {
	now := time.Now()
	todayTime := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second
	endUnixTime := now.Add(-todayTime).Unix()
	startUnixTime := now.Add(time.Duration(-n*24) * time.Hour).Add(-todayTime).Unix()

	var err error
	var result []model.StatsResponse
	resultMap, _ := r.SelectTopIntentByDailyRank(startUnixTime, endUnixTime, n)
	tool.Map2struct(resultMap, &result)

	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	err = middleware.Rdb.Set(ctx, fmt.Sprint("topn_recent", n, "day"), data, 48*time.Hour).Err()
	if err != nil {
		panic(err)
	}
}

// 整理某天的Intent数据，入参dayAgo表示距今已有几天，昨天dayAgo=1
func CollectOneDayIntentCount(daysAgo int) []*model.DailyIntent {
	now := time.Now()
	todayTime := time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second

	startDay := now.Add(time.Duration(-daysAgo*24) * time.Hour).Add(-todayTime)
	endDay := startDay.Add(24 * time.Hour)

	resultMap, err := r.SelectTopIntentIdByTime(startDay.Unix(), endDay.Unix(), 1000)
	if err != nil {
		panic(err)
	}

	var result []model.StatsResponse
	tool.Map2struct(resultMap, &result)

	var dayIntentData []*model.DailyIntent
	for _, record := range result {
		dayIntentData = append(dayIntentData, &model.DailyIntent{UnixTime: startDay.Unix(), IntentId: record.IntentId, Count: record.Cnt})
	}
	err = d.CreateInBatches(dayIntentData, len(dayIntentData))
	if err != nil {
		return nil
	}

	return dayIntentData
}
