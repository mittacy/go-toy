package eMongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var ErrNoInit = errors.New("eMongo: please initialize with eMongo.Init()")

type EMongo struct {
	ConfName   string
	Collection string
}

func (ctl *EMongo) MDB(opts ...*options.CollectionOptions) *mongo.Collection {
	return GetCollection(ctl.ConfName, ctl.Collection, opts...)
}

func GetCollection(name string, collection string, opts ...*options.CollectionOptions) *mongo.Collection {
	if !initFlag {
		panic(ErrNoInit)
	}

	// 获取db
	db, ok := connectPool[name]
	if ok {
		return db.Collection(collection, opts...)
	}

	// 获取配置
	conf, isExist := connectConf[name]
	if !isExist {
		panic("eMongo: " + name + "配置不存在, 请检查配置")
	}

	client, err := GetClient(conf)
	if err != nil {
		l.Errorw("eMongo: 连接失败", "conf", conf, "err", err)
		return &mongo.Collection{}
	}
	db = client.Database(conf.Database)
	connectPool[name] = db

	return db.Collection(collection, opts...)
}

func GetClient(conf Conf) (*mongo.Client, error) {
	if !initFlag {
		panic(ErrNoInit)
	}

	var uri string
	if conf.Password != "" {
		uri = fmt.Sprintf(pswUriFormat, conf.User, conf.Password, conf.Host, conf.Port)
	} else {
		uri = fmt.Sprintf(uriFormat, conf.Host, conf.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// check connect
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
