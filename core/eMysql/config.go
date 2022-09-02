package eMysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dbDSNFormat = "%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local"

// MultipleConf 读写分离配置
type MultipleConf struct {
	Sources  []Conf // 写库
	Replicas []Conf // 读库
}

// ConfToGormDialector 处理配置为gorm配置
// @param conf 配置列表
// @return sourcesDSN 写库dsn列表
// @return replicasDSN 读库dsn列表
func (m *MultipleConf) ConfToGormDialector() (sourcesDialector, replicasDialector []gorm.Dialector) {
	var sources, replicas []gorm.Dialector

	for i := range m.Sources {
		dsn := m.Sources[i].DSN()
		sources = append(sources, mysql.Open(dsn))
	}

	for i := range m.Replicas {
		dsn := m.Replicas[i].DSN()
		replicas = append(replicas, mysql.Open(dsn))
	}

	return sources, replicas
}

type Conf struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	Params   string
}

func (c *Conf) DSN() string {
	dsn := fmt.Sprintf(dbDSNFormat, c.User, c.Password, c.Host, c.Port, c.Database)
	if c.Params != "" {
		dsn = fmt.Sprintf("%s&%s", dsn, c.Params)
	}
	return dsn
}
