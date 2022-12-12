package config

import (
	"fmt"
	"os"
)

var EnvTag = local
var globalConfig = make(map[string]interface{})

const (
	local  = "local"
	docker = "docker"
)

func init() {
	args := os.Args
	for i := range args {
		fmt.Println("args[", i, "]", args[i])
	}

	if len(args) > 1 {
		EnvTag = args[1]
	}

	globalConfig[local] = map[string]string{
		"mysql_url":    "@tcp(localhost:3306)/",
		"mysql_user":   "root",
		"mysql_pass":   "12345678",
		"mysql_dbname": "bot_agent",
		"redis_url":    "localhost:6379",
	}

	globalConfig[docker] = map[string]string{
		"mysql_url":    "@tcp(mysql:3306)/",
		"mysql_user":   "root",
		"mysql_pass":   "12345678",
		"mysql_dbname": "bot_agent",
		"redis_url":    "redis:6379",
	}
}

func GetByEnv(envTag string, key string) string {
	return globalConfig[envTag].(map[string]string)[key]
}

func Get(key string) string {
	return globalConfig[EnvTag].(map[string]string)[key]
}
