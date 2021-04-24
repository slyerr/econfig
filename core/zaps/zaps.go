package zaps

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	// 日志级别
	Level zapcore.Level
	// 日志文件路径
	FilePath string
	// 日志存活时间
	MaxAge time.Duration
	// 日志切割时间单位
	RotationTime time.Duration
}

func NewLogger(config LoggerConfig) *zap.Logger {
	return zap.New(zapcore.NewTee(
		NewFileCore(config),
		NewConsoleCore(config.Level),
	), zap.AddCaller())
}

func NewFileCore(config LoggerConfig) zapcore.Core {
	if len(config.FilePath) == 0 {
		config.FilePath = "./log/%Y-%m-%d.log"
	}
	if config.MaxAge <= 0 {
		config.MaxAge = 30 * 24 * time.Hour
	}
	if config.RotationTime <= 0 {
		config.RotationTime = 24 * time.Hour
	}

	encc := zap.NewProductionEncoderConfig()
	encc.EncodeTime = zapcore.ISO8601TimeEncoder
	encc.EncodeLevel = zapcore.CapitalLevelEncoder
	encc.EncodeCaller = zapcore.FullCallerEncoder
	enc := zapcore.NewConsoleEncoder(encc)

	rotator, err := rotatelogs.New(
		config.FilePath,
		rotatelogs.WithMaxAge(config.MaxAge),
		rotatelogs.WithRotationTime(config.RotationTime))
	if err != nil {
		panic(err)
	}

	ws := zapcore.AddSync(rotator)

	return zapcore.NewCore(enc, ws, config.Level)
}

func NewConsoleCore(level zapcore.Level) zapcore.Core {
	encc := zap.NewProductionEncoderConfig()
	encc.EncodeTime = zapcore.ISO8601TimeEncoder
	encc.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encc.EncodeCaller = zapcore.FullCallerEncoder
	enc := zapcore.NewConsoleEncoder(encc)

	ws := zapcore.AddSync(os.Stdout)

	return zapcore.NewCore(enc, ws, level)
}

func Start(config LoggerConfig) {
	zap.ReplaceGlobals(NewLogger(config))
	zap.S().Warn("zap's global logger has been initialized")
}

func Stop() {
	zap.L().Sync()
}
