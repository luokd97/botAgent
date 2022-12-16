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
	//接收环境变量envTag
	args := os.Args
	for i := range args {
		fmt.Println("args[", i, "]", args[i])
		if args[i] == "-envTag" && len(args) > i+1 {
			EnvTag = args[2]
		}
	}
	fmt.Println("envTag=", EnvTag)

	//配置信息
	globalConfig[local] = map[string]string{
		"mysql_url":    "@tcp(localhost:3306)/",
		"mysql_user":   "root",
		"mysql_pass":   "12345678",
		"mysql_dbname": "bot_agent",
		"redis_url":    "localhost:6379",
		"max_n":        "10000",
	}

	globalConfig[docker] = map[string]string{
		"mysql_url":    "@tcp(mysql:3306)/",
		"mysql_user":   "root",
		"mysql_pass":   "12345678",
		"mysql_dbname": "bot_agent",
		"redis_url":    "redis:6379",
		"max_n":        "10000",
	}
}

func GetByEnv(envTag string, key string) string {
	return globalConfig[envTag].(map[string]string)[key]
}

func Get(key string) string {
	return globalConfig[EnvTag].(map[string]string)[key]
}
