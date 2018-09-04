package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/GO-SQL-Driver/MySQL"
	"fpdxIm/config"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/go-redis/redis"
)

var DB *gorm.DB
var Redis *redis.Client
var Pubsub *redis.PubSub

func Init()  {
	conf := config.Config.DB
	dns := conf["User"].(string) + ":" + conf["Password"].(string) + "@(" + conf["Host"].(string) + ":" + conf["Port"].(string) + ")/" + conf["Name"].(string) + "?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		log.Fatal("数据库连接失败：", err.Error())
	}
	db.DB().SetMaxOpenConns(100)
	DB = db
	rcf := config.Config.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr: rcf["Host"].(string) + ":" + rcf["Port"].(string),
		Password: rcf["Auth"].(string),
		DB: 6,
	})
	if _, err := Redis.Ping().Result(); err != nil {
		log.Fatal("Connect to redis error", err)
		return
	}
	Pubsub = Redis.Subscribe("ws:msgpool:" + config.Config.Address)
	if _, err := Pubsub.Receive(); err != nil {
		panic(err)
	}
}
