package tool

import (
	"encoding/json"
	"fmt"
	"time"
)

// Map2struct s必须传指针类型
func Map2struct(m interface{}, s interface{}) {
	fmt.Printf("map=%v\n", m)
	marshal, err := json.Marshal(m)
	if err != nil {
		fmt.Println("marshal:", err)
		panic(err)
	}
	fmt.Println("after marshal=", string(marshal))
	err = json.Unmarshal(marshal, &s)
	if err != nil {
		fmt.Println("unmarshal:", err)
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("struct=%v\n", s)
}

func TodayUsedTimeDuration() time.Duration {
	nowTime := time.Now()
	return time.Duration(nowTime.Hour())*time.Hour + time.Duration(nowTime.Minute())*time.Minute + time.Duration(nowTime.Second())*time.Second
}

func UnixEpochDaysSince1970(to *time.Time) int {
	if to == nil {
		now := time.Now()
		to = &now
	}
	return int(to.Unix() / int64(24*time.Hour.Seconds()))
}

func LastWeekUnixTimeRange() (start time.Time, end time.Time) {
	nowTime := time.Now()
	todayStartTime := nowTime.Add(-TodayUsedTimeDuration())
	daysUtilLastSunday := nowTime.Weekday()
	if nowTime.Weekday() == time.Sunday {
		daysUtilLastSunday = 7
	}
	start = todayStartTime.Add(time.Duration(-24*(daysUtilLastSunday+7)) * time.Hour)
	end = todayStartTime.Add(time.Duration(-24*daysUtilLastSunday) * time.Hour)
	return
}

func LastMonthUnixTimeRange() (start time.Time, end time.Time) {
	nowTime := time.Now()
	if nowTime.Month() == 1 {
		start = time.Date(nowTime.Year()-1, 12, 1, 0, 0, 0, 0, time.Local)
	} else {
		start = time.Date(nowTime.Year(), nowTime.Month()-1, 1, 0, 0, 0, 0, time.Local)
	}
	end = start.AddDate(0, 1, 0)
	return
}

// 传入一个多返回值的函数调用，只拿第一个返回值
func FirstReturn(fn func(...any) any) any {
	//todo
	return nil
}
