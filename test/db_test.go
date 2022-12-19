package test

import (
	"botApiStats/api"
	"botApiStats/cache"
	"botApiStats/sync"
	"testing"
)

func TestEnvInit(t *testing.T) {
	api.FillTestBotResponse()
	sync.MergeLast3MonthIntentData()
	cache.UpdateAllCache()
}
