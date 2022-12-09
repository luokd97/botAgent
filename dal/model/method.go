package model

import "gorm.io/gen"

// Dynamic SQL
type Method interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)

	// select intent_id, count(*) cnt from `bot_response` group by intent_id order by cnt desc limit @n
	SelectTopIntentId(n int) ([]gen.M, error)

	// select intent_id, count(*) cnt,
	//(select intent_name from bot_response where intent_id=br.intent_id  order by unix_time desc limit 1) as intent_name
	//from `bot_response` br where unix_time between @startTime and @endTime group by intent_id order by cnt desc limit @n
	SelectTopIntentIdByTime(startTime int64, endTime int64, n int) ([]gen.M, error)

	// select intent_name from bot_response where intent_id = @intentId order by created_at desc limit 1
	GetIntentNameByIntentId(intentId string) (string, error)

	// select intent_id, sum(count) total_cnt,
	//(select intent_name from bot_response where intent_id=di.intent_id order by unix_time desc limit 1) as intent_name
	//from `daily_intent` di where unix_time between @startTime and @endTime group by intent_id order by total_cnt desc limit @n
	SelectTopIntentByDailyRank(startTime int64, endTime int64, n int) ([]gen.M, error)

	//select intent_name from (select * from `bot_response` order by unix_time desc limit 100000000) as b
	//where intent_id in (@intentIds)  group by intent_id order by field(intent_id {{for _,id:=range intentIds}},@id{{end}} )
	SelectNewestIntentNamesByIntentIds(intentIds []string) ([]string, error)

	//select intent_id, intent_name from (select * from `bot_response` order by unix_time desc limit 100000000) as b  group by intent_id
	SelectIntentIdMapIntentName() (map[string]string, error)
}
