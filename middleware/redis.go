package middleware

import (
	"github.com/go-redis/redis/v9"
)

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		//Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
