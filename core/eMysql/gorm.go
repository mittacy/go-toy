package eMysql

import (
	"context"
	"errors"
	"fmt"
	stackErrors "github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// GetConnectByConf 获取gorm连接
func GetConnectByConf(conf MultipleConf) (*gorm.DB, error) {
	dsn := conf.Sources[0].DSN()

	db, err := gorm.Open(mysql.Open(dsn), gormConfig())
	if err != nil {
		return nil, err
	}

	if len(conf.Sources) > 1 || len(conf.Replicas) > 0 {
		sources, replicas := conf.ConfToGormDialector()

		dbResolverConfig := dbresolver.Config{
			Sources:  sources,
			Replicas: nil,
			Policy:   dbresolver.RandomPolicy{},
		}

		if len(replicas) > 0 {
			dbResolverConfig.Replicas = replicas
		}

		if err = db.Use(dbresolver.Register(dbResolverConfig)); err != nil {
			return nil, err
		}
	}

	return db, nil
}

var ErrNoInit = errors.New("eMysql: please initialize with eMysql.Init()")

// EGorm 在结构体引入组合并赋值ConfName，即可通过GDB()获取gorm连接
// Example
// type User struct {
// 	 EGorm
// }
//
// var user = User{EGorm{ConfName: "localhost"}}
//
// func (u *User) GetUser(id int64) error {
// 	 u.GDB().Where("id = ?", id).First()
// }
type EGorm struct {
	ConfName string
}

// GDB 获取DB连接
func (ctl *EGorm) GDB() *gorm.DB {
	return GetConnect(ctl.ConfName)
}

// Create 创建
// @param c
// @param values 记录数据
// @return error
func (ctl *EGorm) Create(c context.Context, values interface{}) error {
	if err := ctl.GDB().WithContext(c).Create(values).Error; err != nil {
		return stackErrors.WithStack(err)
	}
	return nil
}

// Save 更新结构体指定id的所有字段
// @param c
// @param records 记录数据
// @param updateColumns 当id为0时，需要更新的字段
// @return error
func (ctl *EGorm) SaveById(c context.Context, records interface{}, updateColumns []string) error {
	if len(updateColumns) == 0 {
		return nil
	}

	if err := ctl.GDB().WithContext(c).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(records).Error; err != nil {
		return err
	}
	return nil
}

const NotLimit = -1

// Updates 更新指定字段
// @param c
// @param table 更新表名
// @param where where条件
// @param noWhere no where条件
// @param updates 更新数据
// @param limit 更新限制条数，不限制则使用NotLimit
// @return int64 更新有效结果数量
// @return error
func (ctl *EGorm) Updates(c context.Context, table string, where, noWhere, updates map[string]interface{}, limit int) (int64, error) {
	if updates == nil || len(updates) == 0 {
		return 0, nil
	}

	dbCtl := ctl.GDB().WithContext(c).Table(table)
	if where != nil && len(where) > 0 {
		dbCtl = dbCtl.Where(where)
	}
	if noWhere != nil && len(noWhere) > 0 {
		dbCtl = dbCtl.Not(noWhere)
	}
	if limit != NotLimit {
		dbCtl = dbCtl.Limit(limit)
	}

	res := dbCtl.Updates(updates)
	if res.Error != nil {
		return res.RowsAffected, stackErrors.WithStack(res.Error)
	}
	return res.RowsAffected, nil
}

// First 查询第一条记录
// @param c
// @param where where条件
// @param noWhere no where条件
// @param result 查询结果数据
// @return error
func (ctl *EGorm) First(c context.Context, where, noWhere map[string]interface{}, result interface{}) error {
	dbCtl := ctl.GDB().WithContext(c)
	if where != nil && len(where) > 0 {
		dbCtl = dbCtl.Where(where)
	}
	if noWhere != nil && len(noWhere) > 0 {
		dbCtl = dbCtl.Not(noWhere)
	}

	if err := dbCtl.First(result).Error; err != nil {
		return stackErrors.WithStack(err)
	}
	return nil
}

func GetConnect(name string) *gorm.DB {
	if !initFlag {
		panic(ErrNoInit)
	}

	if db, ok := connectPool[name]; ok {
		return db
	}

	conf, isExist := connectConf[name]
	if !isExist {
		panic("eGorm: " + name + "配置不存在, 请检查配置")
	}

	db, err := GetConnectByConf(conf)
	if err != nil {
		gormLog.Error(context.Background(), fmt.Sprintf("eGorm: 连接数据库失败, conf: %+v, err: %+v", conf, err))
		return &gorm.DB{}
	}
	connectPool[name] = db
	return db
}

// gormConfig gorm连接配置
func gormConfig() *gorm.Config {
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // 是否禁用表名复数形式
		Logger:         gormLog,
	}
}
