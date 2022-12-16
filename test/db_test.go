package test

import (
	"botApiStats/api"
	"botApiStats/cron"
	"botApiStats/dal"
	"botApiStats/dal/query"
	"testing"
)

var r = query.Use(dal.Db).BotResponse

func TestResponseDataFill(t *testing.T) {
	api.FillTestBotResponse()
}

// 按天整理90天内的bot_response表数据并写入daily_intent
func TestDailyDataFill(t *testing.T) {
	cron.UpdateDailyIntentFrom3MonthAgo()
}

// 整理昨日数据并更新缓存（每日定时任务）
func TestUpdateCache(t *testing.T) {
	cron.DailyUpdateCacheJob()
}

func TestEnvInit(t *testing.T) {
	TestResponseDataFill(nil)
	TestDailyDataFill(nil)
	TestUpdateCache(nil)
}
