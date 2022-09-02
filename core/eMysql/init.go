package eMysql

import "gorm.io/gorm"

var (
	initFlag    bool                    // 初始化标识
	connectConf map[string]MultipleConf // 连接配置
	connectPool map[string]*gorm.DB     // 连接单例池
)

func Init(c map[string]MultipleConf, logOptions []LogConfigOption) {
	// init conf
	connectConf = c

	// init log
	initLog(logOptions...)

	// init pool
	connectPool = make(map[string]*gorm.DB, 0)

	initFlag = true
}
