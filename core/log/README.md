# log
zap框架二次封装

### 使用

```go
func main() {
    // 新建句柄
    field := zapcore.Field{
        Key:    "module_name",
        Type:   zapcore.StringType,
        String: "serverName",
    }
    bizLog := New("default",
        WithPath("./storage/logs"),
        WithTimeFormat("2006-01-02 15:04:05"),
        WithLevel(zapcore.InfoLevel),
        WithPreName("biz_"),
        WithEncoderJSON(false),
        WithFields(field))
    
    // 记录日志
    bizLog.Debug("this is SugarDebug")
    bizLog.Info("this is SugarInfo")
    bizLog.Warn("this is SugarWarn")
    bizLog.Error("this is SugarError")
    
    bizLog.Debugf("this is %s", "Debugf")
    bizLog.Infof("this is %s", "Infof")
    bizLog.Warnf("this is %s", "Warn")
    bizLog.Errorf("this is %s", "Errorf")
    
    bizLog.Debugw("this is Debugw", "k", "Debugw")
    bizLog.Infow("this is Infow", "k", "Infow")
    bizLog.Warnw("this is Warnw", "k", "Warnw")
    bizLog.Errorw("this is Errorw", "k", "Errorw")
}


```

### 全局配置
修改默认配置，之后新建的日志，将默认使用这些配置

SetGlobalConf(options ...ConfigOption), 配置项使用WithXxx()

- WithPath 日志路径，默认为 .
- WithLevel 最低打印级别，默认为 debug
- WithLogInConsole 是否打印到控制台，默认为 false
- WithFields 该日志都加的字段，默认为空
- WithEncoderJSON 是否为json格式日志，默认为 true
- WithTimeFormat 时间格式，默认为 2006-01-02T15:04:05Z07:00
- WithPreName 日志前缀，默认为 biz_

```go
func main() {
    field := zapcore.Field{
        Key:    "module_name",
        Type:   zapcore.StringType,
        String: "helloLog",
    }
    // 设置默认配置后，之后新建的所有日志都默认使用该配置
    SetDefaultConf(
    	WithPath("./storage/logs"),
        WithTimeFormat("2006-01-02 15:04:05"),
        WithLevel(zapcore.InfoLevel),
        WithPreName("biz_"),
        WithEncoderJSON(false),
        WithFields(field))
    
    // 新建日志
    // ……
}
```

### 日志记录上下文信息

```go
func main() {
    traceIdKey := "trace_id"
	c := context.WithValue(context.Background(), traceIdKey, "r61f0ed0d70098_Zw8R1aoyl4tGeB4HMV")

	// 告知log从上下文获取trace_id并记录到每个日志
	SetDefaultConf(WithCtxField(traceIdKey))

	l := New("trace")
	l.DebugWithCtx(c,"this is SugarDebug")
	l.InfoWithCtx(c,"this is SugarInfo")
	l.WarnWithCtx(c,"this is SugarWarn")
	l.ErrorWithCtx(c,"this is SugarError")

	l.DebugfWithCtx(c,"this is %s", "Debugf")
	l.InfofWithCtx(c,"this is %s", "Infof")
	l.WarnfWithCtx(c,"this is %s", "Warn")
	l.ErrorfWithCtx(c,"this is %s", "Errorf")

	l.DebugwWithCtx(c,"this is Debugw", "k", "Debugw")
	l.InfowWithCtx(c,"this is Infow", "k", "Infow")
	l.WarnwWithCtx(c,"this is Warnw", "k", "Warnw")
	l.ErrorwWithCtx(c,"this is Errorw", "k", "Errorw")
}
```