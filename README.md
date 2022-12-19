## 数据库表
#### 1、 bot_response

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `created_at` | 记录创建时间-unix时间戳 | bigint |  | NO |  |  |
| 3 | `agent_id` |  | varchar(128) |  | NO |  |  |
| 4 | `intent_id` |  | varchar(128) |  | NO |  |  |
| 5 | `intent_name` |  | longtext |  | YES |  |  |


#### 2、 daily_intent

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `date` | 记录创建日期-unix时间戳 | bigint |  | NO |  |  |
| 3 | `agent_id` |  | varchar(128) |  | NO |  |  |
| 4 | `intent_id` |  | varchar(128) |  | NO |  |  |
| 5 | `count` |  | bigint |  | NO |  |  |

唯一约束：(date,agent_id,intent_id)同一天只用单条记录来保存某个agent_id相关的intent_id出现次数



## 部署方式
api服务
```sh
$ docker-compose up --build
```
cron定时任务
```sh
$ crontab -e

0 0 2 * * docker-compose --file cron/dock-compose.yml up --build 
```

## 根据接口注释生成swagger.json
Documentation served at http://127.0.0.1:8000/docs
```sh
$ swag init
```

## 根据[generate.go](dal%2Fgenerate%2Fgenerate.go)生成orm代码

```sh
$ go run "dal/generate/generate.go"
```
## 为本地Mysql和Redis导入测试数据
```sh
$ go test ./test -run TestEnvInit
```
