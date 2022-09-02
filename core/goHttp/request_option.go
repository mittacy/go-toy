package goHttp

import (
	"bytes"
	"io"
	"net/http"
	netUrl "net/url"
)

func newRequest(method, url string, options ...RequestOption) (*Request, error) {
	rq := Request{}
	for _, v := range options {
		v(&rq)
	}

	if err := rq.build(method, url); err != nil {
		return nil, err
	}

	return &rq, nil
}

type Request struct {
	request    *http.Request
	paramsUrl  string // 请求传递的url
	requestUrl string // 处理后实际请求的url

	headers    map[string]string
	formParams map[string]string
	jsonData   interface{}
}

func (ro *Request) Request() *http.Request {
	return ro.request
}

func (ro *Request) build(method, url string) error {
	ro.paramsUrl = url
	ro.requestUrl = url

	// 处理form请求表单数据
	if ro.formParams != nil {
		requestUrl, err := netUrl.Parse(url)
		if err != nil {
			return err
		}
		urlValues := netUrl.Values{}
		for k, v := range ro.formParams {
			urlValues.Add(k, v)
		}
		requestUrl.RawQuery = urlValues.Encode()
		ro.requestUrl = requestUrl.String()
	}

	// 处理json请求体
	var body io.Reader
	if ro.jsonData != nil {
		jsonData, err := jsonMarshal(ro.jsonData)
		if err != nil {
			return err
		}
		body = bytes.NewReader(jsonData)
	}

	request, err := http.NewRequest(method, ro.requestUrl, body)
	if err != nil {
		return err
	}

	// 处理头数据
	if ro.headers != nil {
		for k, v := range ro.headers {
			request.Header.Add(k, v)
		}
	}

	// json数据请求
	if ro.jsonData != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	ro.request = request
	return nil
}

type RequestOption func(request *Request)

// 添加请求头
func WithAddHeaders(headers map[string]string) RequestOption {
	return func(request *Request) {
		request.headers = headers
	}
}

// 添加请求表单
func WithFormParams(params map[string]string) RequestOption {
	return func(request *Request) {
		request.formParams = params
	}
}

// 添加JSON请求体
func WithJSONData(data interface{}) RequestOption {
	return func(request *Request) {
		request.jsonData = data
	}
}
