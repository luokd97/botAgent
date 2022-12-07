package cron

import (
	"botApiStats/dal"
	"botApiStats/dal/query"
	"botApiStats/middleware"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

var ctx = context.Background()
var r = query.Use(dal.Db).BotResponse
var datetimeFormat = "2000-01-01 00:00:00.00"

func TestCron() {
	fmt.Println("start cron test")
	c := cron.New()
	EntryID, err := c.AddFunc("@every 5s", func() { fmt.Println("Every 5s") })
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
}

func updateCacheCron() {
	c := cron.New()
	fmt.Println(time.Now(), "每天凌晨2点执行一次")
	EntryID, err := c.AddFunc("0 0 2 * * ?", UpdateCacheJob)
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
}

func UpdateCacheJob() {
	//更新缓存
	fmt.Println(time.Now(), "开始更新缓存")
	updateYesterdayTopN()
	updateLast7dayTopN()
	updateLast30dayTopN()
	fmt.Println(time.Now(), "缓存更新完成")
}

// 更新昨日topN缓存
func updateYesterdayTopN() {
	startTime := time.Now().Add(-24*time.Hour - 2*time.Hour).Format(datetimeFormat)
	endTime := time.Now().Add(-2 * time.Hour).Format(datetimeFormat)

	result, err := r.SelectTopIntentIdByTime(startTime, endTime, 5)
	if err != nil {
		panic(err)
	} else {
		data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		err = middleware.Rdb.Set(ctx, "topn_yesterday", data, 48*time.Hour).Err()
		if err != nil {
			panic(err)
		}
	}
}

// 更新近7天topN缓存
func updateLast7dayTopN() {
	startTime := time.Now().Add(-24*7*time.Hour - 2*time.Hour).Format(datetimeFormat)
	endTime := time.Now().Add(-2 * time.Hour).Format(datetimeFormat)

	result, err := r.SelectTopIntentIdByTime(startTime, endTime, 5)
	if err != nil {
		panic(err)
	} else {
		data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		err = middleware.Rdb.Set(ctx, "topn_last7day", data, 48*time.Hour).Err()
		if err != nil {
			panic(err)
		}
	}
}

// 更新近30天topN缓存
func updateLast30dayTopN() {
	startTime := time.Now().Add(-24*30*time.Hour - 2*time.Hour).Format(datetimeFormat)
	endTime := time.Now().Add(-2 * time.Hour).Format(datetimeFormat)

	result, err := r.SelectTopIntentIdByTime(startTime, endTime, 5)
	if err != nil {
		panic(err)
	} else {
		data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		err = middleware.Rdb.Set(ctx, "topn_last30day", data, 48*time.Hour).Err()
		if err != nil {
			panic(err)
		}
	}
}
