package log

var std *Logger

func Init(options ...ConfigOption) {
	SetDefaultConf(options...)
}

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Sugar = std.Sugar

	Debug = std.Debug
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal

	Debugf = std.Debugf
	Infof = std.Infof
	Warnf = std.Warnf
	Errorf = std.Errorf
	DPanicf = std.DPanicf
	Panicf = std.Panicf
	Fatalf = std.Fatalf

	Debugw = std.Debugw
	Infow = std.Infow
	Warnw = std.Warnw
	Errorw = std.Errorw
	DPanicw = std.DPanicw
	Panicw = std.Panicw
	Fatalw = std.Fatalw

	DebugWithCtx = std.DebugWithCtx
	InfoWithCtx = std.InfoWithCtx
	WarnWithCtx = std.WarnWithCtx
	ErrorWithCtx = std.ErrorWithCtx
	DPanicWithCtx = std.DPanicWithCtx
	PanicWithCtx = std.PanicWithCtx
	FatalWithCtx = std.FatalWithCtx
	DebugfWithCtx = std.DebugfWithCtx
	InfofWithCtx = std.InfofWithCtx
	WarnfWithCtx = std.WarnfWithCtx
	ErrorfWithCtx = std.ErrorfWithCtx
	DPanicfWithCtx = std.DPanicfWithCtx
	PanicfWithCtx = std.PanicfWithCtx
	FatalfWithCtx = std.FatalfWithCtx
	DebugwWithCtx = std.DebugwWithCtx
	InfowWithCtx = std.InfowWithCtx
	WarnwWithCtx = std.WarnwWithCtx
	ErrorwWithCtx = std.ErrorwWithCtx
	DPanicwWithCtx = std.DPanicwWithCtx
	PanicwWithCtx = std.PanicwWithCtx
	FatalwWithCtx = std.FatalwWithCtx
}

func Default() *Logger {
	return std
}

var (
	Sugar = std.Sugar

	Debug  = std.Debug
	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal

	Debugf  = std.Debugf
	Infof   = std.Infof
	Warnf   = std.Warnf
	Errorf  = std.Errorf
	DPanicf = std.DPanicf
	Panicf  = std.Panicf
	Fatalf  = std.Fatalf

	Debugw  = std.Debugw
	Infow   = std.Infow
	Warnw   = std.Warnw
	Errorw  = std.Errorw
	DPanicw = std.DPanicw
	Panicw  = std.Panicw
	Fatalw  = std.Fatalw

	DebugWithCtx   = std.DebugWithCtx
	InfoWithCtx    = std.InfoWithCtx
	WarnWithCtx    = std.WarnWithCtx
	ErrorWithCtx   = std.ErrorWithCtx
	DPanicWithCtx  = std.DPanicWithCtx
	PanicWithCtx   = std.PanicWithCtx
	FatalWithCtx   = std.FatalWithCtx
	DebugfWithCtx  = std.DebugfWithCtx
	InfofWithCtx   = std.InfofWithCtx
	WarnfWithCtx   = std.WarnfWithCtx
	ErrorfWithCtx  = std.ErrorfWithCtx
	DPanicfWithCtx = std.DPanicfWithCtx
	PanicfWithCtx  = std.PanicfWithCtx
	FatalfWithCtx  = std.FatalfWithCtx
	DebugwWithCtx  = std.DebugwWithCtx
	InfowWithCtx   = std.InfowWithCtx
	WarnwWithCtx   = std.WarnwWithCtx
	ErrorwWithCtx  = std.ErrorwWithCtx
	DPanicwWithCtx = std.DPanicwWithCtx
	PanicwWithCtx  = std.PanicwWithCtx
	FatalwWithCtx  = std.FatalwWithCtx
)
