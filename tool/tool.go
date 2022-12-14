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
