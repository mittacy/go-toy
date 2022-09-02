package eMysql

import (
	"github.com/mittacy/go-toy/core/log"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
	"time"
)

type logConf struct {
	Name                 string
	SlowThreshold        time.Duration
	IgnoreRecordNotFound bool
	Level                logger.LogLevel
}

func defaultLogConf() *logConf {
	return &logConf{
		Name:                 "gorm",
		SlowThreshold:        100 * time.Millisecond,
		IgnoreRecordNotFound: false,
		Level:                logger.Info,
	}
}

type LogConfigOption func(conf *logConf)

// WithName 慢日志名，默认为 mysql
func WithName(name string) LogConfigOption {
	return func(conf *logConf) {
		conf.Name = name
	}
}

// WithSlowThreshold 慢日志时间阈值，默认为 100毫秒
func WithSlowThreshold(duration time.Duration) LogConfigOption {
	return func(conf *logConf) {
		conf.SlowThreshold = duration
	}
}

// WithIgnoreRecordNotFound 是否忽略notFound错误，默认为 false
func WithIgnoreRecordNotFound(isIgnore bool) LogConfigOption {
	return func(conf *logConf) {
		conf.IgnoreRecordNotFound = isIgnore
	}
}

// WithLogLevel 是否忽略notFound错误
func WithLogLevel(level logger.LogLevel) LogConfigOption {
	return func(conf *logConf) {
		conf.Level = level
	}
}

// 日志
var gormLog zapgorm2.Logger

func initLog(options ...LogConfigOption) {
	logConf := defaultLogConf()

	for _, option := range options {
		option(logConf)
	}

	gormLog = zapgorm2.New(log.New(logConf.Name).Log())
	gormLog.SlowThreshold = logConf.SlowThreshold
	gormLog.LogLevel = logConf.Level
	gormLog.IgnoreRecordNotFoundError = logConf.IgnoreRecordNotFound
	gormLog.SetAsDefault()
}
