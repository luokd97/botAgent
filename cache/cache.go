package cache

import (
	"botApiStats/config"
	"botApiStats/dal"
	"botApiStats/dal/enum"
	"botApiStats/dal/query"
	"botApiStats/tool"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var d = query.Use(dal.Db).DailyIntent
var maxN, _ = strconv.Atoi(config.Get("max_n"))

func UpdateAllCache() {
	fmt.Println(time.Now(), "开始更新缓存")
	UpdateCacheByDuration(enum.Yesterday)   //缓存昨天topN结果
	UpdateCacheByDuration(enum.Recent7day)  //缓存过去7天topN结果
	UpdateCacheByDuration(enum.Recent30day) //缓存过去30天topN结果
	UpdateCacheByDuration(enum.Recent90day) //缓存过去90天topN结果
	UpdateCacheByDuration(enum.LastWeek)    //缓存上周的topN结果
	UpdateCacheByDuration(enum.LastMonth)   //缓存上月的topN结果
	fmt.Println(time.Now(), "缓存更新完成")
}

func UpdateCacheByDuration(duration enum.QueryDuration) {
	var queryEndDay = tool.TodayUnixDay() - 1
	var queryStartDay int

	switch duration {
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
	result, err := d.SelectTopNIntentByDailyRank(int64(queryStartDay), int64(queryEndDay), maxN, "")

	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	beforeTomorrow := 24*time.Hour - tool.TodayUsedTimeDuration()
	err = dal.Rdb.Set(context.Background(), duration.String(), data, beforeTomorrow).Err()
	fmt.Println(duration.String(), " expiration after:", beforeTomorrow)
	if err != nil {
		panic(err)
	}
}
