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
	fmt.Println("定时任务UpdateCacheJob已启动")
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
	//获取昨日topN列表
	dayIntents := CollectOneDayIntentCount(1)

	//获取map[intent_id] -> 最新intent_name
	idNameMap := make(map[string]string)
	result, _ := r.SelectIntentIdMapIntentName()
	for i := range result {
		in := result[i]
		idNameMap[in["intent_id"].(string)] = in["intent_name"].(string)
	}

	yesterdayTopNIntents := make([]model.IntentResult, 0)
	for _, dayIntent := range dayIntents {
		yesterdayTopNIntents = append(yesterdayTopNIntents, model.IntentResult{
			Count:      dayIntent.Count,
			IntentId:   dayIntent.IntentId,
			IntentName: idNameMap[dayIntent.IntentId],
		})
	}
	data, err := json.Marshal(yesterdayTopNIntents)
	if err != nil {
		panic(err)
	}
	err = middleware.Rdb.Set(ctx, "topn_yesterday", data, 48*time.Hour).Err()
	if err != nil {
		panic(err)
	}

}

// 整理最近n天的intent数据
func updateRecentNDayTopN(n int) {
	nowTime := time.Now()
	startTime := nowTime.Add(time.Duration(-n*24) * time.Hour).Add(-tool.TodayUsedTimeDuration())
	endTime := nowTime.Add(-tool.TodayUsedTimeDuration())

	var err error
	var result []model.IntentResult
	result, _ = r.SelectTopNIntentByDailyRank(startTime.Unix(), endTime.Unix(), n)

	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	err = middleware.Rdb.Set(ctx, fmt.Sprint("topn_recent", n, "day"), data, 48*time.Hour).Err()
	if err != nil {
		panic(err)
	}
}

// 整理某一天的Intent数据，入参dayAgo表示距今已有几天，昨天dayAgo=1
func CollectOneDayIntentCount(daysAgo int) []*model.DailyIntent {
	nowTime := time.Now()
	startDayTime := nowTime.Add(time.Duration(-daysAgo*24) * time.Hour).Add(-tool.TodayUsedTimeDuration())
	endDayTime := startDayTime.AddDate(0, 0, 1)

	var topIntents []model.IntentResult
	topIntents, err := r.SelectTopNIntentByTime(startDayTime.Unix(), endDayTime.Unix(), 1000)
	if err != nil {
		panic(err)
	}

	var dayIntentData []*model.DailyIntent
	for _, intent := range topIntents {
		dayIntentData = append(dayIntentData, &model.DailyIntent{
			Date:     startDayTime.Unix(),
			IntentId: intent.IntentId,
			Count:    intent.Count})
	}
	err = d.CreateInBatches(dayIntentData, len(dayIntentData))
	if err != nil {
		return nil
	}

	return dayIntentData
}
