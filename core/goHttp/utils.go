package goHttp

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"net/http"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func fullUrl(host, uri string) string {
	uri = strings.TrimLeft(uri, "/")
	host = strings.TrimRight(host, "/")
	return host + "/" + uri
}

// 从ctx查询b3追踪数据并注入头部
func injectReqB3Header(c context.Context, key string, request *http.Request) error {
	iSpanCtx := c.Value(key)
	spanCtx, ok := iSpanCtx.(model.SpanContext)
	if ok {
		injectorReq := b3.InjectHTTP(request)
		if err := injectorReq(spanCtx); err != nil {
			return err
		}
	}
	return nil
}
