package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 创建日志记录器
func NewLogger() *zap.Logger {
	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	level := zapcore.DebugLevel

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 创建日志记录器
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger
}

// String 创建字符串字段
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 创建 int64 字段
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Float64 创建 float64 字段
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool 创建布尔字段
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Error 创建错误字段
func Error(err error) zap.Field {
	return zap.Error(err)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Duration 创建时间间隔字段
func Duration(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
