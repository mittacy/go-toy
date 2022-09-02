package bizerr

// 服务错误定义
var (
	Success           = &BizErr{0, "success"}
	Param             = &BizErr{1, "params error"}
	Request           = &BizErr{1, "request error"}
	DebounceIntercept = &BizErr{2, "busy service"}
	Unauthorized      = &BizErr{401, "unauthorized"}
	Forbidden         = &BizErr{403, "forbidden permissions"}
	Unknown           = &BizErr{500, "unknown error"}
)

// 业务错误定义格式: 大模块00:中间模块00:业务模块00
var ()
