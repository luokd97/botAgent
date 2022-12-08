package model

type BotResponse struct {
	ID         uint `gorm:"primarykey"`
	UnixTime   int64
	IntentId   string
	IntentName string
}

type DailyIntent struct {
	ID         uint   `gorm:"primarykey"`
	UnixTime   int64  `json:"unix_time"`
	IntentId   string `json:"intent_id"`
	IntentName string `json:"intent_name"`
	Count      int    `json:"count"`
}

type StatsRequest struct {
	N         int `json:"n"`          //检索数量前n的intent信息，默认为3
	StartTime int `json:"start_time"` //检索范围的起始时间 unix时间戳
	EndTime   int `json:"end_time"`   //检索范围的结束时间 unix时间戳
}

type CommonIntentRequest struct {
	N        int `json:"n"`         //检索数量前n的intent信息，默认为3
	Duration int `json:"duration" ` //duration-查询的时间范围 枚举类型：0.昨天 1.过去7天 2.过去30天 3.过去90天 4.上周汇总 5.上月汇总" enum:"0,1,2,3,4,5
}

// swagger:response StatsResponse
type StatsResponse struct {
	Cnt        int    `json:"cnt"`         //召回次数
	IntentId   string `json:"intent_id"`   //知识点唯一id
	IntentName string `json:"intent_name"` //知识点最新名称
}

func (BotResponse) TableName() string {
	return "bot_response"
}

func (DailyIntent) TableName() string {
	return "daily_intent"
}
