package model

type BotResponse struct {
	ID         uint   `gorm:"primarykey"`
	CreatedAt  int64  `json:"created_at" gorm:"index:idx_created_at;comment:记录创建时间-unix时间戳"`
	IntentId   string `json:"intent_id" gorm:"size:128;index:idx_intent_id;"`
	IntentName string `json:"intent_name"`
}

func (BotResponse) TableName() string {
	return "bot_response"
}

type DailyIntent struct {
	ID       uint   `gorm:"primarykey"`
	Date     int64  `json:"date" gorm:"uniqueIndex:date_intent_uniq;comment:记录创建日期-unix时间戳"`
	IntentId string `json:"intent_id" gorm:"size:128; uniqueIndex:date_intent_uniq"`
	Count    int    `json:"count"`
}

func (DailyIntent) TableName() string {
	return "daily_intent"
}

type IntentResult struct {
	IntentId   string `json:"intent_id"`   //知识点唯一id
	IntentName string `json:"intent_name"` //知识点最新名称
	Count      int    `json:"count"`       //按当前条件统计到的数量
}

// DurationStatsRequest 按枚举范围统计
type DurationStatsRequest struct {
	N        int `json:"n" binding:"required,gte=1,lte=1000"` //检索数量前n的intent信息，n允许范围[1,1000]
	Duration int `json:"duration" binding:"gte=0,lte=5"`      //duration-查询的时间范围 枚举类型：0.昨天 1.过去7天 2.过去30天 3.过去90天 4.上周汇总 5.上月汇总" enum:"0,1,2,3,4,5
}

// ExactStatsRequest 按精确范围统计
type ExactStatsRequest struct {
	N         int   `json:"n" binding:"required,gte=1,lte=1000"`                  //检索数量前n的intent信息，n允许范围[1,1000]
	StartTime int64 `json:"start_time" binding:"required,gte=0"`                  //检索范围的起始时间 unix时间戳
	EndTime   int64 `json:"end_time" binding:"required,gte=0,gtefield=StartTime"` //检索范围的结束时间 unix时间戳
}
