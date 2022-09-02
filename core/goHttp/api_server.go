package goHttp

import (
	"github.com/mittacy/go-toy/core/log"
)

type ApiServer struct {
	Host  string
	reply func() IReply
	log   *log.Logger
}

func NewApiServer(host string, options ...ApiOption) ApiServer {
	s := ApiServer{
		Host:  host,
		reply: DefaultReply,
		log:   gc.log,
	}

	for _, v := range options {
		v(&s)
	}

	return s
}

type ApiOption func(apiServer *ApiServer)

// 自定义日志名
func WithLogName(logName string) ApiOption {
	return func(apiServer *ApiServer) {
		if logName == gc.logName {
			apiServer.log = gc.log
		} else {
			apiServer.log = log.New(logName)
		}
	}
}

// 自定义业务响应结构
// 对 ApiServer.GET、ApiServer.POST 有效
func WithReply(reply IReply) ApiOption {
	return func(apiServer *ApiServer) {
		apiServer.reply = func() IReply {
			return reply
		}
	}
}
