package goHttp

import (
	"github.com/mittacy/go-toy/core/log"
	"net/http"
)

var gc *globalClient

type globalClient struct {
	client  http.Client
	log     *log.Logger
	logName string
	b3Trace
}

func Init(client http.Client, logName string, options ...ClientOption) {
	c := globalClient{
		client:  client,
		logName: logName,
		log:     log.New(logName),
		b3Trace: b3Trace{
			isOpenB3Trace: false,
		},
	}

	for _, v := range options {
		v(&c)
	}

	gc = &c
}

type ClientOption func(client *globalClient)

type b3Trace struct {
	isOpenB3Trace     bool
	b3TraceFromCtxKey string
}

// 是否开启zipkin-B3追踪，将会在头部注入追踪数据，数据从上下文获取
// key为从上下文获取b3.SpanContext数据的键名
func WithHeaderB3TraceFromCtx(key string) ClientOption {
	return func(client *globalClient) {
		client.isOpenB3Trace = true
		client.b3TraceFromCtxKey = key
	}
}
