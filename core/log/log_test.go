package log

import (
	"context"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLog(t *testing.T) {
	Init(WithTimeFormat("2006-01-02 15:04:05"), WithLevel(DebugLevel))

	Debug("this is SugarDebug")
	Info("this is SugarInfo")
	Warn("this is SugarWarn")
	Error("this is SugarError")

	Debugf("this is %s", "Debugf")
	Infof("this is %s", "Infof")
	Warnf("this is %s", "Warn")
	Errorf("this is %s", "Errorf")

	Debugw("this is Debugw", "k", "Debugw")
	Infow("this is Infow", "k", "Infow")
	Warnw("this is Warnw", "k", "Warnw")
	Errorw("this is Errorw", "k", "Errorw")
}

func TestConf(t *testing.T) {
	field := zapcore.Field{
		Key:    "module_name",
		Type:   zapcore.StringType,
		String: "serverName",
	}
	bizLog := New("file",
		WithPath("./storage/logs"),
		WithTimeFormat("2006-01-02 15:04:05"),
		WithLevel(InfoLevel),
		WithPreName("biz_"),
		WithEncoderJSON(false),
		WithFields(field))

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

func TestDefault(t *testing.T) {
	field := zapcore.Field{
		Key:    "module_name",
		Type:   zapcore.StringType,
		String: "serverName",
	}
	SetDefaultConf(WithPath("./storage/logs"),
		WithTimeFormat("2006-01-02 15:04:05"),
		WithLevel(InfoLevel),
		WithPreName("global_"),
		WithEncoderJSON(false),
		WithFields(field))

	bizLog := New("default")

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

func TestWithCtx(t *testing.T) {
	traceIdKey := "trace_id"
	c := context.WithValue(context.Background(), traceIdKey, "r61f0ed0d70098_Zw8R1aoyl4tGeB4HMV")

	// 告知log从上下文获取trace_id并记录到每个日志
	SetDefaultConf(WithCtxField(traceIdKey))

	l := New("trace")
	l.DebugWithCtx(c, "this is SugarDebug")
	l.InfoWithCtx(c, "this is SugarInfo")
	l.WarnWithCtx(c, "this is SugarWarn")
	l.ErrorWithCtx(c, "this is SugarError")

	l.DebugfWithCtx(c, "this is %s", "Debugf")
	l.InfofWithCtx(c, "this is %s", "Infof")
	l.WarnfWithCtx(c, "this is %s", "Warn")
	l.ErrorfWithCtx(c, "this is %s", "Errorf")

	l.DebugwWithCtx(c, "this is Debugw", "k", "Debugw")
	l.InfowWithCtx(c, "this is Infow", "k", "Infow")
	l.WarnwWithCtx(c, "this is Warnw", "k", "Warnw")
	l.ErrorwWithCtx(c, "this is Errorw", "k", "Errorw")

	t.Log("success")
}
