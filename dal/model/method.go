package model

import "gorm.io/gen"

// Dynamic SQL
type Method interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)

	// select intent_id, count(*) cnt from `bot_response` group by intent_id order by cnt desc limit @n
	SelectTopIntentId(n int) ([]gen.M, error)

	// select intent_id, count(*) cnt from `bot_response` where created_at between @startTime and @endTime group by intent_id order by cnt desc limit @n
	SelectTopIntentIdByTime(startTime string, endTime string, n int) ([]gen.M, error)

	// select intent_name from bot_response where intent_id = @intentId order by created_at desc limit 1
	GetIntentNameByIntentId(intentId string) (string, error)
}
