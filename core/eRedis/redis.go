package eRedis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// GetClientByConf 获取redis客户端
// @param conf 配置
// @param defaultDB 默认数据库
// @return *redis.Client
func GetClientByConf(conf Conf, db int) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       db,
	}

	if conf.PoolSize > 0 { // 最大连接数
		options.PoolSize = conf.PoolSize
	}
	if conf.MinIdleConn > 0 { // 最小空闲连接数
		options.MinIdleConns = conf.MinIdleConn
	}
	if conf.IdleTimeout > 0 { // 空闲时间(秒)
		options.IdleTimeout = conf.IdleTimeout * time.Second
	}

	rdb := redis.NewClient(options)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return &redis.Client{}, err
	}

	return rdb, nil
}

var ErrNoInit = errors.New("eRedis: please initialize with eRedis.Init()")

// ERedis 在结构体引入组合并赋值RedisConfName、RedisDB，即可通过RDB()获取redis连接
// Example
// type User struct {
// 	 eRedis.DB
// }
//
// var user = User{ConfName: "localhost", DB: 1}
//
// func (u *User) GetUser() {
// 	 u.RDB().Get(key)
// }
type ERedis struct {
	ConfName string
	DB       int
}

// RDB 获取redis连接
// @return *redis.Client
func (r *ERedis) RDB() *redis.Client {
	return GetClient(r.ConfName, r.DB)
}

func (r *ERedis) Del(c context.Context, keys ...string) error {
	if len(keys) <= 0 {
		return nil
	}
	if err := r.RDB().Del(c, keys...).Err(); err != nil {
		l.ErrorwWithCtx(c, "删除缓存失败", "keys", keys, "err", err)
		return err
	}
	return nil
}

// GetClient 获取redis客户端
// @param name 配置名
// @param defaultDB 默认数据库
// @return *redis.Client
func GetClient(name string, defaultDB int) *redis.Client {
	if !initFlag {
		panic(ErrNoInit)
	}

	cacheName := connectName(name, defaultDB)

	//poolLock.RLock()
	if db, ok := connectPool[cacheName]; ok {
		//poolLock.RUnlock()
		return db
	}
	//poolLock.RUnlock()

	conf, isExist := connectConf[name]
	if !isExist {
		panic("eRedis: " + name + " 配置不存在, 请检查连接配置")
	}

	db, err := GetClientByConf(conf, defaultDB)
	if err != nil {
		l.Errorw("连接数据库失败", "conf", conf, "err", err)
		return &redis.Client{}
	}

	//poolLock.Lock()
	connectPool[cacheName] = db
	//poolLock.Unlock()

	return db
}

func connectName(name string, db int) string {
	return fmt.Sprintf("%s:%d", name, db)
}
