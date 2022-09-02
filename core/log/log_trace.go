package log

import (
	"context"
	"go.uber.org/zap/zapcore"
)

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) DebugWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Debugw(msg, keyValues...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) InfoWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Infow(msg, keyValues...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) WarnWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Warnw(msg, keyValues...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) ErrorWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Errorw(msg, keyValues...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanicWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.DPanicw(msg, keyValues...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) PanicWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Panicw(msg, keyValues...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) FatalWithCtx(c context.Context, msg string, fields ...zapcore.Field) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, fields)
	l.Fatalw(msg, keyValues...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) DebugfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) InfofWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) WarnfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) ErrorfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanicfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) PanicfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) FatalfWithCtx(c context.Context, template string, args ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	l.Sugar().With(keyValues...).Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *Logger) DebugwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) InfowWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) WarnwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) ErrorwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) DPanicwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) PanicwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) FatalwWithCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	keyValues := getCtxAndFieldsKeyVal(c, l.conf.CtxFields, []zapcore.Field{})
	keysAndValues = append(keysAndValues, keyValues...)
	l.Sugar().Fatalw(msg, keysAndValues...)
}

func getCtxAndFieldsKeyVal(c context.Context, ctxFields []string, fields []zapcore.Field) []interface{} {
	keyValues := make([]interface{}, 0, len(fields)+len(ctxFields))
	for _, v := range fields {
		keyValues = append(keyValues, v.Key, v.Interface)
	}
	for _, v := range ctxFields {
		keyValues = append(keyValues, v, c.Value(v))
	}
	return keyValues
}
