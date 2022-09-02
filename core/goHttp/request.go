package goHttp

import (
	"context"
	"errors"
	"io/ioutil"
)

type requestMethod = string

const (
	GetRequest  requestMethod = "GET"
	PostRequest requestMethod = "POST"
)

func (server *ApiServer) check() error {
	// 检查全局配置
	if gc == nil {
		panic("goHttp: 未初始化goHttp, 请调用goHttp.Init()初始化")
	}

	// 检查host
	if server.Host == "" {
		return errors.New("host cannot be empty")
	}

	return nil
}

// GET请求api，解析整个响应数据
// @param uri 请求uri
// @param result 响应结果数据绑定
// @param options 请求可省参数
// @return error 错误
func (server *ApiServer) SimpleGet(c context.Context, uri string, result interface{}, options ...RequestOption) error {
	return server.simple(GetRequest, c, uri, result, options...)
}

// GET请求api，解析整个响应数据
// @param uri 请求uri
// @param params 请求参数
// @param result 响应结果数据绑定
// @return error
func (server *ApiServer) SimpleGetParams(c context.Context, uri string, params map[string]string, result interface{}) error {
	return server.SimpleGet(c, uri, result, WithFormParams(params))
}

// POST请求api，解析整个响应数据
// @param uri 请求uri
// @param result 响应结果数据绑定
// @param options 请求可省参数
// @return error
func (server *ApiServer) SimplePost(c context.Context, uri string, result interface{}, options ...RequestOption) error {
	return server.simple(PostRequest, c, uri, result, options...)
}

// POST请求api，解析整个响应数据
// @param uri 请求uri
// @param body 请求体
// @param result 响应结果数据绑定
// @return error
func (server *ApiServer) SimplePostJSON(c context.Context, uri string, body map[string]interface{}, result interface{}) error {
	return server.SimplePost(c, uri, result, WithJSONData(body))
}

func (server *ApiServer) simple(method requestMethod, c context.Context, uri string, result interface{}, options ...RequestOption) error {
	err := server.check()
	if err != nil {
		return err
	}

	requestUrl := fullUrl(server.Host, uri)
	request, err := newRequest(method, requestUrl, options...)
	if err != nil {
		return err
	}

	if gc.isOpenB3Trace {
		if err := injectReqB3Header(c, gc.b3TraceFromCtxKey, request.Request()); err != nil {
			server.log.ErrorwWithCtx(c, "api请求注入b3追踪头部失败", "uri", uri)
		}
	}

	resp, err := gc.client.Do(request.Request())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		server.log.ErrorwWithCtx(c, requestUrl, "headers", request.Request().Header, "res", string(body), "err", err)
		return err
	}

	server.log.InfowWithCtx(c, request.paramsUrl, "requestUrl", request.requestUrl, "headers", request.headers, "params", request.formParams, "body", request.jsonData, "res", string(body))

	if result != nil {
		if err = jsonUnmarshal(body, result); err != nil {
			return err
		}
	}

	return nil
}

// GET请求api，解析响应数据的data字段数据
// @param uri 请求uri
// @param dataResult 响应结果data数据绑定
// @param options 请求可省参数
// @return int 响应码
// @return error 错误
func (server *ApiServer) Get(c context.Context, uri string, dataResult interface{}, options ...RequestOption) (int, error) {
	return server.doAndParseData(GetRequest, c, uri, dataResult, options...)
}

// GET请求api，解析响应数据的data字段数据
// @param uri 请求uri
// @param params 请求参数
// @param dataResult 响应结果data数据绑定
// @return int 响应码
// @return error 错误
func (server *ApiServer) GetParams(c context.Context, uri string, params map[string]string, dataResult interface{}) (int, error) {
	return server.Get(c, uri, dataResult, WithFormParams(params))
}

// POST请求api，解析响应数据的data字段数据
// @param uri 请求uri
// @param dataResult 响应结果data数据绑定
// @param options 请求可省参数
// @return int 响应码
// @return error 错误
func (server *ApiServer) Post(c context.Context, uri string, dataResult interface{}, options ...RequestOption) (int, error) {
	return server.doAndParseData(PostRequest, c, uri, dataResult, options...)
}

// POST请求api，解析响应数据的data字段数据
// @param uri 请求uri
// @param body 请求体
// @param dataResult 响应结果data数据绑定
// @return int 响应码
// @return error 错误
func (server *ApiServer) PostJSON(c context.Context, uri string, body map[string]interface{}, dataResult interface{}) (int, error) {
	return server.Post(c, uri, dataResult, WithJSONData(body))
}

func (server *ApiServer) doAndParseData(method requestMethod, c context.Context, uri string, dataResult interface{}, options ...RequestOption) (int, error) {
	err := server.check()
	if err != nil {
		return 0, err
	}

	reply := server.reply()

	requestUrl := fullUrl(server.Host, uri)
	request, err := newRequest(method, requestUrl, options...)
	if err != nil {
		return reply.GetUnknownCode(), err
	}

	if gc.isOpenB3Trace {
		if err := injectReqB3Header(c, gc.b3TraceFromCtxKey, request.Request()); err != nil {
			server.log.ErrorwWithCtx(c, "api请求注入b3追踪头部失败", "uri", uri)
		}
	}

	resp, err := gc.client.Do(request.Request())
	if err != nil {
		return reply.GetUnknownCode(), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		server.log.ErrorwWithCtx(c, requestUrl, "headers", request.Request().Header, "res", string(body), "err", err)
		return reply.GetUnknownCode(), err
	}

	server.log.InfowWithCtx(c, request.paramsUrl, "requestUrl", request.requestUrl, "headers", request.headers, "params", request.formParams, "body", request.jsonData, "res", string(body))

	if dataResult == nil {
		return reply.GetSuccessCode(), nil
	}

	if err = jsonUnmarshal(body, reply); err != nil {
		return reply.GetUnknownCode(), err
	}

	if !reply.IsSuccess() {
		return reply.GetCode(), errors.New(reply.GetMsg())
	}

	if err = reply.UnmarshalData(dataResult); err != nil {
		return reply.GetCode(), err
	}

	return reply.GetCode(), nil
}
