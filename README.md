## 数据库表

#### 1、 bot_response

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `unix_time` |  | bigint |  | YES |  |  |
| 3 | `intent_id` |  | longtext |  | YES |  |  |
| 4 | `intent_name` |  | longtext |  | YES |  |  |


#### 2、 daily_intent

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | `id` |  | bigint unsigned | PRI | NO | auto_increment |  |
| 2 | `unix_time` |  | bigint |  | YES |  |  |
| 3 | `intent_id` |  | longtext |  | YES |  |  |
| 4 | `intent_name` |  | longtext |  | YES |  |  |
| 5 | `count` |  | bigint |  | YES |  |  |




## 部署方式

```sh
$ docker-compose up --build
```

## 根据接口注释生成swagger.json
Documentation served at http://127.0.0.1:8000/docs
```sh
$ /doc_generator.sh
```

## 根据[generate.go](cmd%2Fgenerate%2Fgenerate.go)生成orm代码

```sh
$ /doc_generator.sh
```
## 导入测试数据
```sh
$ /doc_generator.sh
```
