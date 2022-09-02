package eMongo

import (
	"github.com/mittacy/go-toy/core/log"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	uriFormat    = "mongodb://%s:%d"       // 无加密
	pswUriFormat = "mongodb://%s:%s@%s:%d" // 加密
)

type Conf struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var (
	l           *log.Logger
	connectConf map[string]Conf            // 连接配置
	initFlag    bool                       // 是否初始化
	connectPool map[string]*mongo.Database // 连接单例池
)

// Init 初始化
// Example:
// c = map[string]Conf{
//		"localhost": {
//			Host:     viper.GetString("MONGO_RW_HOST"),
//			Port:     viper.GetInt("MONGO_RW_PORT"),
//			Database: viper.GetString("MONGO_RW_DATABASE"),
//			User:     viper.GetString("MONGO_RW_USERNAME"),
//			Password: viper.GetString("MONGO_RW_PASSWORD"),
//		},
//	}
// @param c
func Init(c map[string]Conf) {
	connectConf = c
	connectPool = map[string]*mongo.Database{}
	l = log.New("mongo")
	initFlag = true
}
