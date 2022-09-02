package eRedis

import (
	"github.com/go-redis/redis/v8"
	"github.com/mittacy/go-toy/core/log"
	"time"
)

type Conf struct {
	Host        string
	Password    string
	Port        int
	PoolSize    int
	MinIdleConn int
	IdleTimeout time.Duration
}

var (
	initFlag    bool                     // 初始化标识
	connectConf map[string]Conf          // 连接配置
	connectPool map[string]*redis.Client // 连接单例池
	//poolLock    sync.RWMutex             // 单例池锁	// Deprecated 加快访问速度移除锁，重复创建也可接受
	l *log.Logger // redis日志
)

// Init 初始化
// Example:
// connectConf = map[string]Conf{
//		"localhost": {
//			Host:        viper.GetString("REDIS_LOCALHOST_RW_HOST"),
//			Password:    viper.GetString("REDIS_LOCALHOST_RW_PASSWORD"),
//			Port:        viper.GetInt("REDIS_LOCALHOST_RW_PORT"),
//			PoolSize:    viper.GetInt("REDIS_LOCALHOST_POOL_SIZE"),
//			MinIdleConn: viper.GetInt("REDIS_LOCALHOST_MIN_IDLE_CONN"),
//			IdleTimeout: viper.GetDuration("REDIS_LOCALHOST_IDLE_TIMEOUT"),
//		},
//	}
func Init(connectConf map[string]Conf) {
	// 初始化redis连接配置池
	initConf(connectConf)

	// 初始化单例池
	connectPool = make(map[string]*redis.Client, 0)
	//poolLock = sync.RWMutex{}

	// 初始化日志
	l = log.New("redis")

	initFlag = true
}

func initConf(c map[string]Conf) {
	connectConf = c
}
