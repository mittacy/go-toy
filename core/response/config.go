package response

var (
	ctxFieldKeys []string
)

// 增加自定义响应字段，从上下文获取
func WithFieldKeysFromCtx(keys ...string) {
	ctxFieldKeys = append(ctxFieldKeys, keys...)
}
