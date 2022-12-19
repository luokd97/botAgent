package main

import (
	"botApiStats/cache"
	"botApiStats/sync"
)

func main() {
	//整理昨日bot_response数据到daily_intent表
	sync.MergeOneDayIntentData(1)
	//更新业务缓存
	cache.UpdateAllCache()
}
