package model

// Dynamic SQL
type Method interface {

	// select intent_id, agent_id, count(*) count,
	//(select intent_name from bot_response where intent_id=br.intent_id  order by created_at desc limit 1) as intent_name
	//from `bot_response` br where {{if agentId != ""}}agent_id=@agentId and {{end}}
	//created_at between @startTime and @endTime group by intent_id order by count desc limit @n
	SelectTopNIntentByTime(startTime int64, endTime int64, n int, agentId string) ([]IntentResult, error)

	// select intent_id,agent_id, sum(count) count,
	//(select intent_name from bot_response where intent_id=di.intent_id order by created_at desc limit 1) as intent_name
	//from `daily_intent` di where {{if agentId != ""}}agent_id=@agentId and {{end}}
	//date between @startTime and @endTime group by intent_id order by count desc limit @n
	SelectTopNIntentByDailyRank(startTime int64, endTime int64, n int, agentId string) ([]IntentResult, error)

	//select intent_id, intent_name from (select * from `bot_response` order by created_at desc limit 100000000) as b  group by intent_id
	SelectIntentIdMapIntentName() (idNameMap []map[string]interface{}, err error)
}
