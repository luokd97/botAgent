## 数据库表
#### 1、 bot_response

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `created_at` | 记录创建时间-unix时间戳 | bigint |  | NO |  |  |
| 3 | `intent_id` |  | varchar(128) |  | NO |  |  |
| 4 | `intent_name` |  | longtext |  | YES |  |  |


#### 2、 daily_intent

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `date` | 记录创建日期-unix时间戳 | bigint |  | NO |  |  |
| 3 | `intent_id` |  | varchar(128) |  | NO |  |  |
| 4 | `count` |  | bigint |  | NO |  |  |

唯一约束：(date,intent_id)同一天只用单条记录来保存某个intent_id出现次数



## 部署方式
本地构建，容器内运行
```sh
$ go build -o app .
$ docker-compose up --build
```

## 根据接口注释生成swagger.json
Documentation served at http://127.0.0.1:8000/docs
```sh
$ swag init
```

## 根据[generate.go](cmd%2Fgenerate%2Fgenerate.go)生成orm代码

```sh
$ go run "cmd/generate/generate.go"
```
## 为本地Mysql和Redis导入测试数据
```sh
$ go test ./test -run TestEnvInit
```
