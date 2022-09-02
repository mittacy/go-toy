package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

type Logger struct {
	name string
	l    *zap.Logger
	conf Conf
}

func (l *Logger) Log() *zap.Logger {
	return l.l
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.l.Sugar()
}

var logPool = map[string]*Logger{}

// New 创建新日志文件句柄，使用默认配置
// @param name 日志名
// @param options 日志配置，将覆盖默认配置
// @return *Logger
func New(name string, options ...ConfigOption) *Logger {
	if strings.Trim(name, " ") == "" {
		panic("日志名不能为空")
	}

	// 检查共用同名日志句柄
	if l, ok := logPool[name]; ok && l != nil {
		return l
	}

	// 创建新日志
	c := getDefaultConf()
	for _, option := range options {
		option(&c)
	}
	c.Name = name

	l := newWithConf(c)
	logPool[name] = l
	return l
}

func newWithConf(conf Conf) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "log_at",
		LevelKey:    "level",
		NameKey:     "logger",
		MessageKey:  "context",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder, // 大写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(conf.TimeFormat))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 是否为JSON格式日志
	var encoder zapcore.Encoder
	if conf.IsJSONEncode {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 是否输出控制台
	var writerSyncs []zapcore.WriteSyncer
	if conf.LogInConsole {
		writerSyncs = append(writerSyncs, zapcore.AddSync(os.Stdout))
	}

	var cores []zapcore.Core
	// 区别日志级别
	if conf.IsDistinguish {
		infoFileName := conf.LogPath() + "_info"
		errFileName := conf.LogPath() + "_error"

		cores = []zapcore.Core{
			zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(append(writerSyncs, zapcore.AddSync(getWriter(infoFileName)))...),
				zap.LevelEnablerFunc(func(level zapcore.Level) bool {
					return level >= conf.LowLevel && level < zapcore.ErrorLevel
				}),
			),
			zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(append(writerSyncs, zapcore.AddSync(getWriter(errFileName)))...),
				zap.LevelEnablerFunc(func(level zapcore.Level) bool {
					return level >= zap.ErrorLevel
				}),
			),
		}
	} else {
		cores = []zapcore.Core{
			zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(append(writerSyncs, zapcore.AddSync(getWriter(conf.LogPath())))...),
				conf.LowLevel,
			),
		}
	}

	// 全局字段添加到每个日志中
	l := &Logger{
		name: conf.Name,
		l:    zap.New(zapcore.NewTee(cores...), zap.Fields(conf.Fields...)),
		conf: conf,
	}

	return l
}

func getWriter(fileName string) io.Writer {
	path := fileName + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// 不存在创建文件
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if file, err = os.Create(path); err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}

	return file
}
