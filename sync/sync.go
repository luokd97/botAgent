package sync

import (
	"botApiStats/config"
	"botApiStats/dal"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"botApiStats/tool"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var d = query.Use(dal.Db).DailyIntent
var r = query.Use(dal.Db).BotResponse
var maxN, _ = strconv.Atoi(config.Get("max_n"))

// MergeOneDayIntentData 整理某一天的Intent数据写入daily_intent表，入参dayAgo表示距今已有几天，昨天dayAgo=1
func MergeOneDayIntentData(daysAgo int) []*model.DailyIntent {
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
			Date:     int64(tool.TodayUnixDay() - daysAgo),
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

// MergeLast3MonthIntentData 手动触发，整理三个月内的response更新到daily_intent表
func MergeLast3MonthIntentData() {
	for i := 1; i <= 90; i++ {
		MergeOneDayIntentData(i)
	}
}
