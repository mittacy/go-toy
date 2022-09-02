package response

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mittacy/go-toy/core/bizerr"
	"github.com/mittacy/go-toy/core/log"
	"net/http"
)

var NilData = struct{}{}

func Success(c *gin.Context, data interface{}) {
	Custom(c, http.StatusOK, bizerr.Success.Code, bizerr.Success.Msg, data)
}

func SuccessMsg(c *gin.Context, data interface{}, msg string) {
	Custom(c, http.StatusOK, bizerr.Success.Code, msg, data)
}

func SuccessNil(c *gin.Context) {
	Custom(c, http.StatusOK, bizerr.Success.Code, bizerr.Success.Msg, NilData)
}

func FailMsg(c *gin.Context, msg string) {
	Custom(c, http.StatusOK, bizerr.Request.Code, msg, NilData)
}

func FailErr(c *gin.Context, err error) {
	CustomErr(c, http.StatusOK, bizerr.Code(err), err, NilData)
}

func Unknown(c *gin.Context, err error) {
	Custom(c, http.StatusOK, bizerr.Code(err), bizerr.Unknown.Error(), struct{}{})
}

// Unauthorized 未认证
func Unauthorized(c *gin.Context) {
	Custom(c, http.StatusOK, bizerr.Unauthorized.Code, bizerr.Unauthorized.Error(), NilData)
}

// Forbidden 权限不足
func Forbidden(c *gin.Context) {
	Custom(c, http.StatusOK, bizerr.Forbidden.Code, bizerr.Forbidden.Error(), NilData)
}

func Custom(c *gin.Context, httpCode, bizCode int, msg string, data interface{}) {
	CustomErr(c, httpCode, bizCode, errors.New(msg), data)
}

func CustomErr(c *gin.Context, httpCode, bizCode int, err error, data interface{}) {
	msg := err.Error()
	if bizCode == bizerr.Unknown.Code {
		msg = bizerr.Unknown.Msg
	}

	h := map[string]interface{}{
		"code": bizCode,
		"msg":  msg,
		"data": data,
	}
	for _, v := range ctxFieldKeys {
		h[v] = c.MustGet(v)
	}

	if bizCode != bizerr.Success.Code && gin.Mode() != gin.ReleaseMode && !bizerr.IsBizErr(err) {
		h["err_stack"] = fmt.Sprintf("%+v", err)
	}

	c.JSON(httpCode, h)
}

// ValidateErr 表单解析错误响应
func ValidateErr(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		log.ErrorwWithCtx(c, "request params err", "err", err.Error())
		// 非validator错误
		Unknown(c, err)
		return
	}

	// validator错误进行翻译
	details := removeTopStruct(errs.Translate(trans))

	// 随机返回校验错误中的一条到 msg 字符串
	msg := "param error"
	for _, v := range details {
		msg = v
		break
	}

	log.ErrorwWithCtx(c, "request params err", "err", msg)
	Custom(c, http.StatusOK, bizerr.Param.Code, msg, details)
	return
}

// FailCheckBizErr 检查错误是否为业务错误，否则记录日志并响应未知
func FailCheckBizErr(c *gin.Context, log *log.Logger, req interface{}, title string, err error) {
	if !bizerr.IsBizErr(err) {
		log.ErrorwWithCtx(c, title, "req", req, "err", err)
	}

	CustomErr(c, http.StatusOK, bizerr.Code(err), err, NilData)
	return
}
