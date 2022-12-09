package middleware

import (
	"botApiStats/config"
	"github.com/go-redis/redis/v9"
)

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Get("redis_url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
