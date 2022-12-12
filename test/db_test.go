package test

import (
	"botApiStats/cron"
	"botApiStats/dal"
	"botApiStats/dal/model"
	"botApiStats/dal/query"
	"fmt"
	"testing"
)

var r = query.Use(dal.Db).BotResponse

// 写入机器人api调用模拟数据
func TestResponseDataFill(t *testing.T) {
	var simpleResponseData = []*model.BotResponse{
		{UnixTime: 1664899200, IntentId: "0", IntentName: "intent0"}, //22-10-5
		{UnixTime: 1668787200, IntentId: "1", IntentName: "intent1"}, //22-11-19
		{UnixTime: 1669305600, IntentId: "2", IntentName: "intent2"}, //22-11-25
		{UnixTime: 1669478400, IntentId: "1", IntentName: "intent1"}, //22-11-27
		{UnixTime: 1669651200, IntentId: "0", IntentName: "intent0"}, //22-11-29
		{UnixTime: 1669737600, IntentId: "2", IntentName: "intent2"}, //22-11-30
		{UnixTime: 1669824000, IntentId: "3", IntentName: "intent3"}, //22-12-1
		{UnixTime: 1669910400, IntentId: "2", IntentName: "intent2"}, //22-12-2
		{UnixTime: 1669910400, IntentId: "2", IntentName: "intent2"}, //22-12-2
		{UnixTime: 1669997000, IntentId: "3", IntentName: "intent3"}, //22-12-3
		{UnixTime: 1670457600, IntentId: "3", IntentName: "intent3"}, //22-12-8
	}
	err := r.CreateInBatches(simpleResponseData, len(simpleResponseData))
	if err != nil {
		panic(err)
		return
	}
}

// 按天整理60天内的bot_response表数据并写入daily_intent
func TestDailyDataFill(t *testing.T) {
	for i := 1; i <= 60; i++ {
		cron.CollectOneDayIntentCount(i)
	}
}

// 手动生成检索缓存（正式环境通过定时任务生成）
func TestUpdateCache(t *testing.T) {
	cron.UpdateCacheJob()
}

func TestEnvInit(t *testing.T) {
	TestResponseDataFill(nil)
	TestDailyDataFill(nil)
	TestUpdateCache(nil)
}

func TestSelectNewestIntentNamesByIntentIds(t *testing.T) {
	names, err := r.SelectNewestIntentNamesByIntentIds([]string{"1", "2"})
	if err != nil {
		panic(err)
	}
	fmt.Println(names)
}

func TestTrue(t *testing.T) {
	fmt.Println("true")
}
